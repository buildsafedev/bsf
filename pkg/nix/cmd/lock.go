package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Lock generates the Nix flake lock file
func Lock() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd := exec.Command("nix", "flake", "lock", fmt.Sprintf("path:%s/bsf/", dir))

	// Connect the command's stdin, stdout, and stderr to the terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed with %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 130 { // ctrl+c
				return nil
			}
		}
		return fmt.Errorf("failed with %s", err)
	}

	return nil
}
