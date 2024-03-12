package generate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	golang "github.com/buildsafedev/bsf/pkg/generate/golang"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	btemplate "github.com/buildsafedev/bsf/pkg/nix/template"
)

var supportedLanguages = []string{"Golang"}

// Generate reads bsf.hcl, resolves dependencies and generates bsf.lock, bsf/flake.nix, bsf/default.nix, etc.
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

	lang := "go"
	if conf.GoModule != nil {
		lang = "go"
	}

	cr := hcl2nix.ResolveCategoryRevisions(conf.Packages, lockPackages)
	err = btemplate.GenerateFlake(btemplate.Flake{
		// Description:         "bsf flake",
		NixPackageRevisions: cr.Revisions,
		DevPackages:         cr.Development,
		RuntimePackages:     cr.Runtime,
		Language:            lang,
	}, fh.FlakeFile)
	if err != nil {
		return err
	}

	err = GenAppModule(fh, conf)
	if err != nil {
		return err
	}

	return nil
}

// GenAppModule will generate default.nix file or other files necessary to build the app based on programming language
func GenAppModule(fh *hcl2nix.FileHandlers, conf *hcl2nix.Config) error {
	if conf.GoModule != nil {
		goMod2NixPath := filepath.Join("bsf/", "gomod2nix.toml")
		outFile := goMod2NixPath
		pkgs, err := golang.GenGolangPackages("./", goMod2NixPath, 10)
		if err != nil {
			return fmt.Errorf("error generating pkgs: %v", err)
		}

		var goPackagePath string
		var subPackages []string

		output, err := golang.Marshal(pkgs, goPackagePath, subPackages)
		if err != nil {
			return fmt.Errorf("error marshaling output: %v", err)
		}

		err = os.WriteFile(outFile, output, 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %v", err)
		}
		err = btemplate.GenerateGoModule(conf.GoModule, fh.DefFlakeFile)
		if err != nil {
			return err
		}
	}

	return nil

}
