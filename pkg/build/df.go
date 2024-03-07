package build

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"golang.org/x/sync/errgroup"
)

var (
	dockerFileTmpl = `# Nix builder
FROM nixos/nix:latest AS builder

# Copy our source and setup our working dir.
COPY . /tmp/build
WORKDIR /tmp/build/bsf

# Build runtime package dependencies
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build

# Build additional packages we need for runtime
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build .#runtimeEnvs.{{ .Platform }}.runtime -o runtimeEnv

{{ if (.DevDeps)}}
# Build development packages if devDeps is set to true in bsf.hcl
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    --option filter-syscalls false \
    build .#devEnvs.{{ .Platform }}.development -o devEnv
{{ end }}

# Copy the Nix store closure into a directory. The Nix store closure is the
# entire set of Nix store values that we need for our build and custom environment.
RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure
RUN cp -R $(nix-store -qR runtimeEnv/) /tmp/nix-store-closure
{{ if (.DevDeps)}}
RUN cp -R $(nix-store -qR devEnv/) /tmp/nix-store-closure
{{ end }}

# # Final image is based on scratch. We copy a bunch of Nix dependencies
# # but they're fully self-contained so we don't need Nix anymore.
{{ if (.DevDeps)}}
FROM busybox
{{ else }}
FROM scratch
{{ end }}

WORKDIR /result
{{ if ne .Config ""}}
COPY {{ .Config }} /result/app
{{ end }}
# Copy /nix/store
COPY --from=builder /tmp/nix-store-closure /nix/store
# Add symlink to result
COPY --from=builder /tmp/build/bsf/result /result
COPY --from=builder /tmp/build/bsf/runtimeEnv /result/env
{{ if (.DevDeps)}}
COPY --from=builder /tmp/build/bsf/devEnv /result/env
{{ end }}
# Add /result/env to the PATH
ENV SSL_CERT_FILE="/result/env/etc/ssl/certs/ca-bundle.crt"
ENV PATH="/result/env/bin:${PATH}"
{{ if gt (len .EnvVars) 0 }}ENV {{ range $key, $value := .EnvVars }}{{ $key }}={{ quote $value }} {{ end }}{{ end }}
{{ if gt (len .Cmd) 0 }}CMD [{{ range $index, $element := .Cmd }} {{if $index}}, {{end}} "{{ quote $element }}" {{ end }}]{{ end }}
{{ if gt (len .Entrypoint) 0 }} ENTRYPOINT [{{ range $index, $element := .Entrypoint }}{{if $index}}, {{end}} "{{ quote $element }}" {{ end }}]{{ end }}
`
)

type buildOpts struct {
	DockerFileLoc string
	Name          string
}

func dockerbuild(ctx context.Context, opts buildOpts) error {
	// todo: we should introduce Builder config in bsf.hcl to allow users to configure custom builders for Nix/Builkit
	c, err := autoClient(ctx)
	if err != nil {
		return err
	}

	pipeR, pipeW := io.Pipe()
	solveOpts, err := newSolveOpt(opts, pipeW)
	if err != nil {
		return err
	}
	ch := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := c.Solve(ctx, nil, *solveOpts, ch)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		d, err := progressui.NewDisplay(os.Stderr, progressui.TtyMode)
		if err != nil {
			// If an error occurs while attempting to create the tty display,
			// fallback to using plain mode on stdout (in contrast to stderr).
			d, _ = progressui.NewDisplay(os.Stdout, progressui.PlainMode)
		}
		// not using shared context to not disrupt display but let is finish reporting errors
		_, err = d.UpdateFrom(context.TODO(), ch)
		return err
	})

	// for some reason the error message for loading dockerfile is hidden if we have the below routine.
	// TODO: we need a proper fix for this.
	if os.Getenv("BSF_DEBUG_DOCKER_BUILD") != "1" {
		eg.Go(func() error {
			if err := dockerLoad(ctx, pipeR, pipeW); err != nil {
				return err
			}
			return pipeR.Close()
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func newSolveOpt(opts buildOpts, w io.WriteCloser) (*client.SolveOpt, error) {
	buildCtx := "."

	file := opts.DockerFileLoc

	localDirs := map[string]string{
		"context":    buildCtx,
		"dockerfile": filepath.Dir(file),
	}

	// todo: use gateway
	frontend := "dockerfile.v0"

	// if clicontext.Bool("clientside-frontend") {
	// 	frontend = ""
	// }

	frontendAttrs := map[string]string{
		"filename": filepath.Base(file),
	}

	// for _, buildArg := range clicontext.StringSlice("build-arg") {
	// 	kv := strings.SplitN(buildArg, "=", 2)
	// 	if len(kv) != 2 {
	// 		return nil, errors.Errorf("invalid build-arg value %s", buildArg)
	// 	}
	// 	frontendAttrs["build-arg:"+kv[0]] = kv[1]
	// }
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: "docker", // TODO: use containerd image store when it is integrated to Docker
				Attrs: map[string]string{
					"name": opts.Name,
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		LocalDirs:     localDirs,
		Frontend:      frontend,
		FrontendAttrs: frontendAttrs,
	}, nil
}

func dockerLoad(ctx context.Context, pipeR io.Reader, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "docker", "load")
	cmd.Stdin = pipeR

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// write both stdout and stderr to w
	go io.Copy(w, stdout)
	go io.Copy(w, stderr)

	return cmd.Run()
}
