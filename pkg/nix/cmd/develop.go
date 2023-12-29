package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// Develop opens a BSF development shell
func Develop() error {
	cmd := exec.Command("nix", "develop", "bsf/.#devShell")

	// Connect the command's stdin, stdout, and stderr to the terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
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
