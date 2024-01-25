package build

import (
	"context"
	"html/template"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

type dockerfileCfg struct {
	Platform   string
	Cmd        []string
	Entrypoint []string
}

// Build builds the environment
func Build(env hcl2nix.ExportConfig) error {
	tmpDir := os.TempDir()
	fh, err := os.Create(tmpDir + "/" + generateRandomFilename())
	if err != nil {
		return err
	}
	defer fh.Close()

	err = GenerateDockerfile(fh, env)
	if err != nil {
		return err
	}

	err = dockerbuild(context.Background(), buildOpts{
		DockerFileLoc: fh.Name(),
		Name:          env.Name,
	})

	return err
}

// GenerateDockerfile generates Dockerfile
func GenerateDockerfile(w io.Writer, env hcl2nix.ExportConfig) error {
	dfc := convertExportCfgToDockerfileCfg(env)

	dftmpl, err := template.New("Dockerfile").Parse(dockerFileTmpl)
	if err != nil {
		return err
	}

	err = dftmpl.Execute(w, dfc)
	if err != nil {
		return err
	}

	return nil
}

func convertExportCfgToDockerfileCfg(env hcl2nix.ExportConfig) dockerfileCfg {
	switch env.Platform {
	case "linux/amd64":
		env.Platform = "x86_64-linux"
	case "linux/arm64":
		env.Platform = "aarch64-linux"
	}

	return dockerfileCfg{
		Platform:   env.Platform,
		Cmd:        env.Cmd,
		Entrypoint: env.Entrypoint,
	}
}

// generateRandomFilename generates a random filename
func generateRandomFilename() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
