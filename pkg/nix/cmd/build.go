package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/config"
	"github.com/buildsafedev/bsf/pkg/generate"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// Build invokes nix build to build the project
func Build(conf *config.Config) error {
	fh, err := hcl2nix.NewFileHandlers(true)
	if err != nil {
		return err
	}

	defer fh.ModFile.Close()
	defer fh.LockFile.Close()
	defer fh.FlakeFile.Close()
	defer fh.DefFlakeFile.Close()

	sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
	if err != nil {
		return fmt.Errorf("error creating search client: %s", err.Error())
	}

	err = generate.Generate(fh, sc)
	if err != nil {
		return err
	}

	cmd := exec.Command("nix", "build", "bsf/.")

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
