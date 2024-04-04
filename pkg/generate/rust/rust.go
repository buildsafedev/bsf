package generate

import (
	"os/exec"
)

func GenCargoNix() error {
	// Run the command
	cmd := exec.Command(
		"nix", "run", "github:cargo2nix/cargo2nix",
	)
	cmd.Dir = "bsf/"
	// Execute the command
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
