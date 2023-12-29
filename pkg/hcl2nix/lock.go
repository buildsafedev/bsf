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

// CategoryRevision holds category revision map  and revision list
type CategoryRevision struct {
	Development map[string]string
	Runtime     map[string]string
	Revisions   []string
}

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

	if !strings.Contains(pkg, "@") {
		// NOTE- we require user to explicitly mention package versions.
		// This allows us to regenerate lock file and other Nix files from bsf.hcl itself.
		// Go also has similar experience where user has to explicitly mention package versions in go.mod and the CLI can resolve "@latest" for imperative UX.
		return nil, fmt.Errorf("Version not specified for package %s", pkg)
	}
	s := strings.Split(pkg, "@")
	name := s[0]
	version := s[1]

	desiredVersion, err = sc.GetPackageVersion(ctx, name, version)
	if err != nil {
		return nil, fmt.Errorf("error fetching package %s@%s: %v", name, version, err)
	}

	return desiredVersion, nil
}

// ResolveLatestPackageVersion resolves latest package version
func ResolveLatestPackageVersion(ctx context.Context, sc *search.Client, pkg string) (*search.Package, error) {
	if strings.Contains(pkg, "@") {
		return nil, fmt.Errorf("Version is specified for the package %s", pkg)
	}

	versionList, err := sc.ListPackageVersions(ctx, pkg)
	if err != nil {
		return nil, fmt.Errorf("error fetching package %s: %v", pkg, err)
	}

	if len(versionList) == 0 {
		return nil, fmt.Errorf("no versions found for package %s", pkg)
	}

	return &versionList[0], nil
}

// ResolveCategoryRevisions maps packages to their category, returns development packages, runtime packages and a list of Nix revisions
func ResolveCategoryRevisions(pkgs Packages, pkgVersions []search.Package) *CategoryRevision {
	devRevisions := make(map[string]string, 0)
	rtRevisions := make(map[string]string, 0)

	devMap := sliceToMap(pkgs.Development)
	rtMap := sliceToMap(pkgs.Runtime)

	var revisions []string
	for _, pkg := range pkgVersions {
		revisions = append(revisions, pkg.Revision)
		pkgName := getPkgName(pkg.Name)
		if _, ok := devMap[pkgName]; ok {
			devRevisions[pkgName] = pkg.Revision
		}
		if _, ok := rtMap[pkgName]; ok {
			rtRevisions[pkgName] = pkg.Revision
		}
	}
	return &CategoryRevision{
		Development: devRevisions,
		Runtime:     rtRevisions,
		Revisions:   bstrings.SliceToSet(revisions),
	}
}

func sliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool, 0)
	for _, s := range slice {
		m[getPkgName(s)] = true
	}

	return m
}

// getPkgName returns package name without version(if one is present)
func getPkgName(pkg string) string {
	if !strings.Contains(pkg, "@") {
		return pkg
	}

	s := strings.Split(pkg, "@")
	return s[0]
}
