package generate

import (
	"os"
	"os/exec"
)

func GenCargoNix(outFile string) error {
	// Run the command
	cmd := exec.Command(
		"nix", "run", "github:cargo2nix/cargo2nix",
	)
	outputFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	cmd.Stdout = outputFile
	cmd.Stderr = outputFile

	// Execute the command
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
