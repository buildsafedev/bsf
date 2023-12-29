package generate

import (
	"context"
	"os"
	"time"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	btemplate "github.com/buildsafedev/bsf/pkg/nix/template"
)

// Generate reads bsf.hcl, resolves dependencies and generates bsf.lock, bsf/flake.nix and bsf/default.nix
func Generate(fh *hcl2nix.FileHandlers, sc *search.Client) error {
	data, err := os.ReadFile("bsf.hcl")
	if err != nil {
		return err
	}

	conf, err := hcl2nix.ReadConfig(data)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	allPackages, err := hcl2nix.ResolvePackages(ctx, sc, conf.Packages)
	if err != nil {
		return err
	}

	err = hcl2nix.GenerateLockFile(allPackages, fh.LockFile)
	if err != nil {
		return err
	}

	cr := hcl2nix.ResolveCategoryRevisions(conf.Packages, allPackages)
	err = btemplate.GenerateFlake(btemplate.Flake{
		// Description:         "bsf flake",
		NixPackageRevisions: cr.Revisions,
		DevPackages:         cr.Development,
		RuntimePackages:     cr.Runtime,
	}, fh.FlakeFile)
	if err != nil {
		return err
	}
	// todo: there should be a generic method "GenerateApplicationModule" that can switch between different project types
	err = btemplate.GenerateGoModule(conf.GoModule, fh.DefFlakeFile)
	if err != nil {
		return err
	}

	return nil
}
