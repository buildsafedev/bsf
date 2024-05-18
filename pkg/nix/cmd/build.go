package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Build invokes nix build to build the project
func Build(dir string, attribute string) error {
	if attribute == "" {
		attribute = "bsf/."
	}
	cmd := exec.Command("nix", "build", attribute, "-o", dir)

	cmd.Stdout = os.Stdout
	// TODO: in future- we can pipe to stderr pipe and modify error messages to be understandable by the user
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}
	return nil
}
