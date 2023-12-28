package init

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

func TestMapPackageCategory(t *testing.T) {
	tests := []struct {
		name        string
		pkgs        hcl2nix.Packages
		pkgVersions []search.Package
		devExpected map[string]string
		rtExpected  map[string]string
		revExpected []string
	}{
		{
			name: "test case 1",
			pkgs: hcl2nix.Packages{
				Development: []string{"pkg1", "pkg3", "pkg5"},
				Runtime:     []string{"pkg2", "pkg4", "pkg5"},
			},
			pkgVersions: []search.Package{
				{Name: "pkg1", Revision: "rev1", Version: "1.0.0"},
				{Name: "pkg2", Revision: "rev2", Version: "1.1.0"},
				{Name: "pkg3", Revision: "rev3", Version: "1.2.0"},
				{Name: "pkg4", Revision: "rev4", Version: "23.11.0"},
				{Name: "pkg5", Revision: "rev1", Version: "1.0.0"},
			},
			devExpected: map[string]string{
				"pkg1": "rev1",
				"pkg3": "rev3",
				"pkg5": "rev1",
			},
			rtExpected: map[string]string{
				"pkg2": "rev2",
				"pkg4": "rev4",
				"pkg5": "rev1",
			},
			revExpected: []string{"rev1", "rev2", "rev3", "rev4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devGot, rtGot, revGot := mapPackageCategory(tt.pkgs, tt.pkgVersions)
			if !cmp.Equal(devGot, tt.devExpected) {
				t.Errorf("DevDeps got %v, want %v", devGot, tt.devExpected)
			}
			if !cmp.Equal(rtGot, tt.rtExpected) {
				t.Errorf(" RTDeps got %v, want %v", rtGot, tt.rtExpected)
			}
			less := func(a, b string) bool { return a < b }
			equalIgnoreOrder := cmp.Diff(revGot, tt.revExpected, cmpopts.SortSlices(less)) == ""
			if !equalIgnoreOrder {
				t.Errorf("Revisions got %v, want %v", revGot, tt.revExpected)
			}
		})
	}
}
