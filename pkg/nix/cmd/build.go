package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Build invokes nix build to build the project
func Build() (string, error) {
	cmd := exec.Command("nix", "build", "bsf/.")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd.Run() failed with %s", cmd.Stderr)
	}

	return stdout.String(), nil
}
