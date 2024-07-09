package builddocker

import (
	"bufio"
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

// ModifyDockerfile modifies the Dockerfile with the specified tag
func ModifyDockerfile(path, tag string, dev bool) error {
	var dockerfilePath string
	if path != "" {
		dockerfilePath = path + "/Dockerfile"
	} else {
		dockerfilePath = "./Dockerfile"
	}

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile not found")
	}

	file, err := os.Open(dockerfilePath)
	if err != nil {
		return fmt.Errorf("error opening Dockerfile: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading Dockerfile: %v", err)
	}

	var searchTag string
	if dev {
		searchTag = "# bsfimage:dev"
	} else {
		searchTag = "# bsfimage:runtime"
	}

	var selectedFrom string
	var selectedIndex int
	for i, line := range lines {
		if strings.Contains(line, searchTag) {
			selectedFrom = line
			selectedIndex = i
			break
		}
	}

	if selectedFrom == "" {
		return fmt.Errorf("no FROM command found with tag %s", searchTag)
	}

	fromParts := strings.Fields(selectedFrom)
	if len(fromParts) < 2 {
		return fmt.Errorf("invalid FROM command format")
	}

	var newFrom string
	if strings.Contains(fromParts[1], ":") {
		imageParts := strings.Split(fromParts[1], ":")
		newFrom = fmt.Sprintf("FROM %s:%s", imageParts[0], tag)
	} else {
		newFrom = fmt.Sprintf("FROM %s:%s", fromParts[1], tag)
	}
	for _, part := range fromParts[2:] {
		newFrom = fmt.Sprintf("%s %s", newFrom, part)
	}

	lines[selectedIndex] = newFrom

	err = os.WriteFile(dockerfilePath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("error writing to Dockerfile: %v", err)
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
		DevDeps:    env.DevDeps,
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
