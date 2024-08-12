package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Develop opens a BSF development shell
func Develop(pureshell bool) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	var shells []string
	shell := os.Getenv("SHELL")

	// using the `path:` will let users work on the project without having to interact with git
	if pureshell {
		// Prepare the command with `--ignore-environment` and necessary `--keep` variables
		keepVars := []string{
			"HOME", "USER", "LOGNAME", "DISPLAY", "TERM", "IN_NIX_SHELL",
			"NIX_SHELL_PRESERVE_PROMPT", "TZ", "PAGER", "NIX_BUILD_SHELL", "SHLVL",
			"NIX_PATH", "NIX_SSL_CERT_FILE", "NIX_SSL_CERT_DIR",
		}
		// Add --keep flags for each variable
		var keepFlags []string
		for _, v := range keepVars {
			keepFlags = append(keepFlags, "--keep", v)
		}

		args := append([]string{"nix", "develop", "--ignore-environment"}, keepFlags...)
		args = append(args, fmt.Sprintf("path:%s/bsf/.#devShell", dir))

		cmd = exec.Command(args[0], args[1:]...)
	} else if shell == "" {
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
