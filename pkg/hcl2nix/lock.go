package hcl2nix

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"sort"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	bstrings "github.com/buildsafedev/bsf/pkg/strings"
	"github.com/buildsafedev/bsf/pkg/update"
)

// Category defines category of package
type Category int

const (
	// Development is a category of package
	Development Category = iota
	// Runtime is a category of package
	Runtime = 1
)

// CategoryRevision holds category revision map  and revision list
type CategoryRevision struct {
	Development map[string]string
	Runtime     map[string]string
	Revisions   []string
}

// LockFile represents strcuture for LockFile
type LockFile struct {
	App      LockApp       `json:"app"`
	Packages []LockPackage `json:"packages"`
}

// LockApp represents a app
type LockApp struct {
	Name string `json:"name"`
}

// LockPackage represents a package
type LockPackage struct {
	Package *buildsafev1.Package `json:"package"`
	Runtime bool                 `json:"runtime"`
}

// GenerateLockFile generates lock file
func GenerateLockFile(conf *Config, packages []LockPackage, wr io.Writer) error {
	la := LockApp{}

	// In future, when we have more languages, we can check all of them and pick the one that is used.
	if conf.GoModule != nil {
		la.Name = conf.GoModule.Name
	}

	if conf.RustApp != nil {
		la.Name = conf.RustApp.CrateName
	}

	if conf.PoetryApp != nil {
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		pc, err := parsePyProject(currentDir + "/" + conf.PoetryApp.Pyproject)
		if err != nil {
			return err
		}

		la.Name = pc.Tool.Poetry.Name
	}

	if conf.JsNpmApp != nil {
		la.Name = conf.JsNpmApp.PackageName
	}

	lf := LockFile{
		App:      la,
		Packages: packages,
	}

	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}

	if _, err := wr.Write(data); err != nil {
		return err
	}

	return nil
}

// ResolvePackages resolves a list of packages concurrently
func ResolvePackages(ctx context.Context, sc buildsafev1.SearchServiceClient, packages Packages, pkgType string) ([]LockPackage, error) {
	var selectedPackages []string
	switch pkgType {
	case "runtime":
		selectedPackages = packages.Runtime
	default:
		selectedPackages = append(packages.Development, packages.Runtime...)
	}

	allPackages := slices.Compact(selectedPackages)
	resolvedPackages := make([]LockPackage, 0, len(allPackages))
	pkgMap := mapPackageCategory(packages)

	errStr := ""
	var wg sync.WaitGroup
	for _, pkg := range allPackages {
		wg.Add(1)
		go func(pkg string) {
			defer wg.Done()
			p, err := resolvePackage(ctx, sc, pkg)
			if err != nil {
				errStr += fmt.Sprintf("error resolving package %s: %v\n", pkg, err)
				return
			}
			if p.Name == "" {
				errStr += fmt.Sprintf("error resolving package %s: no package found\n", pkg)
				return
			}

			categories := pkgMap[p.Name+"@"+p.Version]
			found := false
			for _, cat := range categories {
				if cat == Runtime {
					found = true
				}
			}

			lp := LockPackage{
				Package: p,
				Runtime: found,
			}
			resolvedPackages = append(resolvedPackages, lp)
		}(pkg)
	}
	wg.Wait()
	if errStr != "" {
		return nil, fmt.Errorf(errStr)
	}
	sort.Slice(resolvedPackages, func(i, j int) bool {
		pi, pj := resolvedPackages[i].Package, resolvedPackages[j].Package
		if pi.Name != pj.Name {
			return pi.Name < pj.Name
		}
		return pi.Version < pj.Version
	})

	return resolvedPackages, nil
}
// ResolvePackage resolves package name
func resolvePackage(ctx context.Context, sc buildsafev1.SearchServiceClient, pkg string) (*buildsafev1.Package, error) {
	var desiredVersion *buildsafev1.FetchPackageVersionResponse
	var err error

	if !strings.Contains(pkg, "@") {
		// NOTE- we require user to explicitly mention package versions.
		// This allows us to regenerate lock file and other Nix files from bsf.hcl itself.
		// Go also has similar experience where user has to explicitly mention package versions in go.mod and the CLI can resolve "@latest" for imperative UX.
		return nil, fmt.Errorf("Version not specified for package %s", pkg)
	}
	name, version := update.TrimVersionInfo(pkg)

	desiredVersion, err = sc.FetchPackageVersion(ctx, &buildsafev1.FetchPackageVersionRequest{
		Name:    name,
		Version: version,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching package %s@%s: %v", name, version, err)
	}

	return desiredVersion.Package, nil
}

// ResolveCategoryRevisions maps packages to their category, returns development packages, runtime packages and a list of Nix revisions
func ResolveCategoryRevisions(pkgs Packages, pkgVersions []LockPackage) *CategoryRevision {
	devRevisions := make(map[string]string, 0)
	rtRevisions := make(map[string]string, 0)

	pkgsMap := mapPackageCategory(pkgs)

	var revisions []string
	for _, pkg := range pkgVersions {
		revisions = append(revisions, pkg.Package.Revision)

		nameWithVersion := pkg.Package.Name + "@" + pkg.Package.Version
		categories := pkgsMap[nameWithVersion]
		if categories == nil {
			continue
		}

		for _, cat := range categories {
			if cat == Runtime {
				// Preference should be given to Attribute name when it is available.
				// Over time, we expect the all packages to have one but to avoid rebuilding the search index, we fall back to name.
				name := pkg.Package.AttrName
				if name == "" {
					name = pkg.Package.Name
				}
				rtRevisions[name] = pkg.Package.Revision
			}

			if cat == Development {
				name := pkg.Package.AttrName
				if name == "" {
					name = pkg.Package.Name
				}
				devRevisions[name] = pkg.Package.Revision
			}
		}

	}

	return &CategoryRevision{
		Development: devRevisions,
		Runtime:     rtRevisions,
		Revisions:   bstrings.SliceToSet(revisions),
	}
}

func mapPackageCategory(packages Packages) map[string][]Category {
	m := make(map[string][]Category)

	addCategory := func(p []string, category Category) {
		for _, pkg := range p {
			name, version := update.TrimVersionInfo(pkg)
			key := name + "@" + version

			m[key] = append(m[key], category)
		}
	}

	addCategory(packages.Runtime, Runtime)
	addCategory(packages.Development, Development)

	return m
}
