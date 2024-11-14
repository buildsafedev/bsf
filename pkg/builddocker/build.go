package builddocker

import (
	"fmt"
	"html/template"
	"io"
	"os/exec"
	"strings"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

type dockerfileCfg struct {
	Platform   string
	Cmd        []string
	Entrypoint []string
	EnvVars    map[string]string
	Config     string
}

func quote(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

// GenerateDockerfile generates Dockerfile
func GenerateDockerfile(w io.Writer, env hcl2nix.OCIArtifact, platform string) error {
	dfc := convertExportCfgToDockerfileCfg(env, platform)

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

func convertExportCfgToDockerfileCfg(env hcl2nix.OCIArtifact, platform string) dockerfileCfg {
	switch platform {
	case "linux/amd64":
		platform = "x86_64-linux"
	case "linux/arm64":
		platform = "aarch64-linux"
	}
	envVarsMap := convertEnvsToMap(env.EnvVars)

	return dockerfileCfg{
		Platform:   platform,
		Cmd:        env.Cmd,
		Entrypoint: env.Entrypoint,
		EnvVars:    envVarsMap,
	}
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
