package build

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

type dockerfileCfg struct {
	Platform   string
	Cmd        []string
	Entrypoint []string
	EnvVars    map[string]string
	DevDeps    bool
	Config     string
}

// Build builds the environment
func Build(env hcl2nix.ExportConfig) error {
	tmpDir, err := createTempDir()
	if err != nil {
		return err
	}
	fh, err := os.Create(tmpDir + "/" + "Dockerfile" + env.Environment + "." + generateRandomFilename())
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

func quote(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

// GenerateDockerfile generates Dockerfile
func GenerateDockerfile(w io.Writer, env hcl2nix.ExportConfig) error {
	dfc := convertExportCfgToDockerfileCfg(env)

	dftmpl, err := template.New("Dockerfile").Funcs(template.FuncMap{
		"quote": quote,
	}).
		Parse(dockerFileTmpl)
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
	envVarsMap := convertEnvsToMap(env.EnvVars)

	return dockerfileCfg{
		Platform:   env.Platform,
		Cmd:        env.Cmd,
		Entrypoint: env.Entrypoint,
		EnvVars:    envVarsMap,
		DevDeps:    env.DevDeps,
		Config:     env.Config,
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

func createTempDir() (string, error) {
	tmpDir := os.TempDir()
	bsfDir := filepath.Join(tmpDir, "bsf")

	if _, err := os.Stat(bsfDir); os.IsNotExist(err) {
		err := os.Mkdir(bsfDir, 0755)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return bsfDir, nil
}

func convertEnvsToMap(envs []string) map[string]string {
	envMap := make(map[string]string)

	for _, env := range envs {
		keyValuePair := strings.SplitN(env, "=", 2)
		if len(keyValuePair) == 2 {
			envMap[keyValuePair[0]] = keyValuePair[1]
		}
	}

	return envMap
}

// GetSnapshotter gets the containerd snapshotter value
func GetSnapshotter() (string, error) {
	script := exec.Command("docker", "info", "-f", " '{{ .DriverStatus }}' ")
	out, err := script.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error fetching  DriverStatus: %s", err)
	}
	return string(out), nil
}
