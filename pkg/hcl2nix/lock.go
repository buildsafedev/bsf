package hcl2nix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	bstrings "github.com/buildsafedev/bsf/pkg/strings"
)

// GenerateLockFile generates lock file
func GenerateLockFile(packages []search.Package, wr io.Writer) error {
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}

	if _, err := wr.Write(data); err != nil {
		return err
	}

	return nil
}

// ResolvePackages resolves a list of packages concurrently
func ResolvePackages(ctx context.Context, sc *search.Client, packages Packages) ([]search.Package, error) {
	allPackages := bstrings.SliceToSet(append(packages.Development, packages.Runtime...))
	resolvedPackages := make([]search.Package, 0, len(allPackages))

	errStr := ""
	var wg sync.WaitGroup
	for _, pkg := range allPackages {
		wg.Add(1)
		go func(pkg string) {
			defer wg.Done()
			p, err := ResolvePackage(ctx, sc, pkg)
			if err != nil {
				errStr += fmt.Sprintf("error resolving package %s: %v\n", pkg, err)
				return
			}

			resolvedPackages = append(resolvedPackages, *p)
		}(pkg)
	}
	wg.Wait()
	if errStr != "" {
		return nil, fmt.Errorf(errStr)
	}

	return resolvedPackages, nil
}

// ResolvePackage resolves package name
func ResolvePackage(ctx context.Context, sc *search.Client, pkg string) (*search.Package, error) {
	var desiredVersion *search.Package
	var err error

	if strings.Contains(pkg, "@") {
		s := strings.Split(pkg, "@")
		name := s[0]
		version := s[1]

		desiredVersion, err = sc.GetPackageVersion(ctx, name, version)
		if err != nil {
			return nil, fmt.Errorf("error fetching package %s@%s: %v", name, version, err)
		}

	} else {
		versions, err := sc.ListPackageVersions(ctx, pkg)
		if err != nil {
			return nil, fmt.Errorf("package %s not found", pkg)
		}

		if versions == nil || len(versions) == 0 {
			return nil, fmt.Errorf("package %s not found", pkg)
		}

		sortedVersions := search.SortPackagesWithTimestamp(versions)
		desiredVersion = &sortedVersions[0]
	}

	return desiredVersion, nil
}
