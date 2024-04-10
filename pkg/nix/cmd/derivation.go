package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GetDrvPathFromResult returns the derivation
func GetDrvPathFromResult(output string, symlink string) (string, error) {
	// TODO: check how to do this via go-nix package-
	// found that it this information comes from narinfo but couldn't figure out how to get narinfo from go-nix
	cmd := exec.Command("nix-store", "--query", "--deriver", output+ symlink)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed with %s", cmd.Stderr)
	}

	return strings.TrimSuffix(stdout.String(), "\n"), nil
}
