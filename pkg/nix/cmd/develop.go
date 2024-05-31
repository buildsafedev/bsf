package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Develop opens a BSF development shell
func Develop() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	var shells []string
	shell := os.Getenv("SHELL")

	// using the `path:` will let users work on the project without having to interact with git
	if shell == "" {
		cmd = exec.Command("nix", "develop", fmt.Sprintf("path:%s/bsf/.#devShell", dir))
	} else {
		if strings.Contains(shell, " ") {
			shells = strings.Split(shell, " ")
		} else if strings.Contains(shell, ":") {
			shells = strings.Split(shell, ":")
		} else {
			shells = strings.Split(shell, ";")
		}

		cmd = exec.Command("nix", "develop", fmt.Sprintf("path:%s/bsf/.#devShell", dir), "-c", shells[0])

	}
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
