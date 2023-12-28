package init

import (
	"strings"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
	bstrings "github.com/buildsafedev/bsf/pkg/strings"
)

func generatehcl2NixConf(pt langdetect.ProjectType) hcl2nix.Config {
	switch pt {
	case langdetect.GoModule:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{
				Development: []string{"go@1.21.5", "gotools@0.7.0", "delve@1.21.2"},
			},
		}
	default:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{},
		}
	}
}

// mapPackageCategory maps packages to their category, returns development packages, runtime packages and a list of Nix revisions
func mapPackageCategory(pkgs hcl2nix.Packages, pkgVersions []search.Package) (map[string]string, map[string]string, []string) {
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

	return devRevisions, rtRevisions, bstrings.SliceToSet(revisions)
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
