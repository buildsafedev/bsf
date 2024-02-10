package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// NixVersion returns the current nix version
func NixVersion() (string, error) {
	var nixVersion string

	script := exec.Command("nix", "--version")
	out, err := script.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error fetching nix version: %s", err)
	}

	splitOut := strings.Fields(string(out))
	if len(splitOut) > 2 {
		nixVersion = splitOut[2]
	} else {
		return "", fmt.Errorf("could not determine the version")
	}

	nixVersion = ("v" + nixVersion)

	return nixVersion, nil
}

// NixShowConfig returns a map of nix configuration values
func NixShowConfig() (map[string]string, error) {

	script := exec.Command("nix", "show-config")
	out, err := script.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error fetching nix config: %s", err)
	}

	config := parseNixConfig(string(out))

	return config, nil
}

func parseNixConfig(config string) map[string]string {
	nixConfig := make(map[string]string)
	for _, line := range strings.Split(config, "\n") {
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			nixConfig[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return nixConfig
}
