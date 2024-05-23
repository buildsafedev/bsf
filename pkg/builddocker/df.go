package builddocker

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

WORKDIR /bin
{{ if ne .Config ""}}
COPY {{ .Config }} /result/app
{{ end }}
# Copy /nix/store
COPY --from=builder /tmp/nix-store-closure /nix/store
# Add symlink to result
COPY --from=builder /tmp/build/bsf/result/bin /bin
COPY --from=builder /tmp/build/bsf/runtimeEnv /bin
{{ if (.DevDeps)}}
COPY --from=builder /tmp/build/bsf/devEnv /bin
{{ end }}
ENV SSL_CERT_FILE="/bin/etc/ssl/certs/ca-bundle.crt"
ENV PATH="/bin:${PATH}"
{{ if gt (len .EnvVars) 0 }}ENV {{ range $key, $value := .EnvVars }}{{ $key }}={{ quote $value }} {{ end }}{{ end }}
{{ if gt (len .Cmd) 0 }}CMD [{{ range $index, $element := .Cmd }} {{if $index}}, {{end}} "{{ quote $element }}" {{ end }}]{{ end }}
{{ if gt (len .Entrypoint) 0 }} ENTRYPOINT [{{ range $index, $element := .Entrypoint }}{{if $index}}, {{end}} "{{ quote $element }}" {{ end }}]{{ end }}
`
)
