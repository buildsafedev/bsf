package generate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	btemplate "github.com/buildsafedev/bsf/pkg/nix/template"
)

// Generate reads bsf.hcl, resolves dependencies and generates bsf.lock, bsf/flake.nix and bsf/default.nix
func Generate(fh *hcl2nix.FileHandlers, sc buildsafev1.SearchServiceClient) error {
	data, err := os.ReadFile("bsf.hcl")
	if err != nil {
		return err
	}

	var dstErr bytes.Buffer
	conf, err := hcl2nix.ReadConfig(data, &dstErr)
	if err != nil {
		return fmt.Errorf("%v", &dstErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	lockPackages, err := hcl2nix.ResolvePackages(ctx, sc, conf.Packages)
	if err != nil {
		return err
	}

	err = hcl2nix.GenerateLockFile(conf, lockPackages, fh.LockFile)
	if err != nil {
		return err
	}

	cr := hcl2nix.ResolveCategoryRevisions(conf.Packages, lockPackages)
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
	if conf.GoModule == nil {
		return nil
	}
	err = btemplate.GenerateGoModule(conf.GoModule, fh.DefFlakeFile)
	if err != nil {
		return err
	}

	return nil
}
