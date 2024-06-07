package builddocker

// source for this file- https://github.com/project-copacetic/copacetic/blob/main/pkg/buildkit/drivers.go
import (
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
)

// GetCurrentContext gets the current context
func GetCurrentContext() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {

		return "", err
	}

	dockerConfigPath := filepath.Join(homeDir, ".docker", "config.json")

	data, err := os.ReadFile(dockerConfigPath)
	if err != nil {
		return "", err
	}

	// Get the current context from the Docker config file
	currentContext := gjson.GetBytes(data, "currentContext").String()
	if currentContext == "" {
		return "", err
	}

	return currentContext, nil

}

// ReadContextEndpoints reads the Docker context endpoints
func ReadContextEndpoints() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(homeDir + "/.docker/contexts/meta")
	if err != nil {
		return nil, err
	}

	endpoints := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		data, err := os.ReadFile(homeDir + "/.docker/contexts/meta/" + file.Name() + "/meta.json")
		if err != nil {
			return nil, err
		}
		contextName := gjson.GetBytes(data, "Name").String()
		endpoints[contextName] = gjson.GetBytes(data, "Endpoints.docker.Host").String()
	}

	return endpoints, nil
}
