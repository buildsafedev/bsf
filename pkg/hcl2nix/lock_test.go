package hcl2nix

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

func TestResolveCategoryRevisions(t *testing.T) {
	tests := []struct {
		name        string
		pkgs        Packages
		pkgVersions []LockPackage
		devExpected map[string]string
		rtExpected  map[string]string
		revExpected []string
	}{
		{
			name: "test case 1",
			pkgs: Packages{
				Development: []string{"pkg1@v1.0.0", "pkg3@1.2.0", "pkg5@1.0.0"},
				Runtime:     []string{"pkg2@1.1.0", "pkg4@23.11.0", "pkg5@1.0.0"},
			},
			pkgVersions: []LockPackage{
				{
					Package: &buildsafev1.Package{
						Name:     "pkg1",
						Revision: "rev1",
						Version:  "1.0.0",
					},
					Runtime: false,
				},
				{
					Package: &buildsafev1.Package{
						Name:     "pkg2",
						Revision: "rev2",
						Version:  "1.1.0",
					},
					Runtime: true,
				},
				{
					Package: &buildsafev1.Package{
						Name:     "pkg3",
						Revision: "rev3",
						Version:  "1.2.0",
					},
					Runtime: false,
				},
				{
					Package: &buildsafev1.Package{
						Name:     "pkg4",
						Revision: "rev4",
						Version:  "23.11.0",
					},
					Runtime: true,
				},
				{
					Package: &buildsafev1.Package{
						Name:     "pkg5",
						Revision: "rev1",
						Version:  "1.0.0",
					},
					Runtime: true,
				},
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
			ct := ResolveCategoryRevisions(tt.pkgs, tt.pkgVersions)
			if !cmp.Equal(ct.Development, tt.devExpected) {
				t.Errorf("DevDeps got %v, want %v", ct.Development, tt.devExpected)
			}
			if !cmp.Equal(ct.Runtime, tt.rtExpected) {
				t.Errorf("RTDeps got %v, want %v", ct.Runtime, tt.rtExpected)
			}
			less := func(a, b string) bool { return a < b }
			equalIgnoreOrder := cmp.Diff(ct.Revisions, tt.revExpected, cmpopts.SortSlices(less)) == ""
			if !equalIgnoreOrder {
				t.Errorf("Revisions got %v, want %v", ct.Revisions, tt.revExpected)
			}
		})
	}
}

func TestMapPackageCategory(t *testing.T) {
	tests := []struct {
		name     string
		packages Packages
		want     map[string][]Category
	}{
		{
			name: "Test Case 1",
			packages: Packages{
				Runtime:     []string{"go@1.20"},
				Development: []string{"node@14.0"},
			},
			want: map[string][]Category{
				"go@1.20":   {Runtime},
				"node@14.0": {Development},
			},
		},
		{
			name: "Test Case 2",
			packages: Packages{
				Runtime:     []string{"go@1.20"},
				Development: []string{"go@1.20", "node@14.0"},
			},
			want: map[string][]Category{
				"go@1.20":   {Runtime, Development},
				"node@14.0": {Development},
			},
		},

		{
			name: "Test Case 2",
			packages: Packages{
				Runtime:     []string{"go@1.20", "node@14.0"},
				Development: []string{"go@1.20", "node@14.0"},
			},
			want: map[string][]Category{
				"go@1.20":   {Runtime, Development},
				"node@14.0": {Runtime, Development},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapPackageCategory(tt.packages); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapPackageCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
