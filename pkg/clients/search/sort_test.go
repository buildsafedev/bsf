package search

import (
	"testing"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/google/go-cmp/cmp"
)

func TestSortPackages(t *testing.T) {
	tests := []struct {
		name string
		pkgs []*buildsafev1.Package
		want []*buildsafev1.Package
	}{
		{
			name: "Test case for semver compliant packages",
			pkgs: []*buildsafev1.Package{
				{
					Name:         "semver",
					Version:      "1.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "semver",
					Version:      "2.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {Name: "semver",
					Version:      "1.5.6",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "0.3.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "11.6.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "2.11.6",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "2.6.9",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
			},
			want: []*buildsafev1.Package{
				{
					Name:         "semver",
					Version:      "11.6.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "2.11.6",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "2.6.9",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "2.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "1.5.6",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "1.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "semver",
					Version:      "0.3.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
			},
		},
		{
			name: "Test case semver and non-semver packages",
			pkgs: []*buildsafev1.Package{
				{
					Name:         "non-semver",
					Version:      "234ca.b243.cc32c",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "non-semver",
					Version:      "3213.122a2.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "semver",
					Version:      "1.5.22",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "semver",
					Version:      "2.6.11",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "non-semver",
					Version:      "213f.4353.75v4",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 5,
				}, {
					Name:         "non-semver",
					Version:      "313f.4353.75v4",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 4,
				}, {
					Name:         "semver",
					Version:      "4.74.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
			},

			want: []*buildsafev1.Package{
				{
					Name:         "semver",
					Version:      "4.74.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "semver",
					Version:      "2.6.11",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "semver",
					Version:      "1.5.22",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				}, {
					Name:         "non-semver",
					Version:      "213f.4353.75v4",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 5,
				}, {
					Name:         "non-semver",
					Version:      "313f.4353.75v4",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 4,
				},
				{
					Name:         "non-semver",
					Version:      "3213.122a2.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
				{
					Name:         "non-semver",
					Version:      "234ca.b243.cc32c",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
			},
		}, {
			name: " Test case SortPackagesWithTimestamp",
			pkgs: []*buildsafev1.Package{
				{
					Name:         "non-semver",
					Version:      "32fd.12a12.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 10,
				}, {
					Name:         "non-semver",
					Version:      "4fd2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 4,
				}, {
					Name:         "non-semver",
					Version:      "5fd2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				}, {
					Name:         "non-semver",
					Version:      "232e.5v33.743",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 7,
				}, {
					Name:         "non-semver",
					Version:      "23r.2324.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 6,
				}, {
					Name:         "non-semver",
					Version:      "6343.4r32.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 3,
				}, {
					Name:         "non-semver",
					Version:      "21d2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 12,
				},
			},
			want: []*buildsafev1.Package{
				{
					Name:         "non-semver",
					Version:      "21d2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 12,
				}, {
					Name:         "non-semver",
					Version:      "32fd.12a12.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 10,
				}, {
					Name:         "non-semver",
					Version:      "232e.5v33.743",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 7,
				}, {
					Name:         "non-semver",
					Version:      "23r.2324.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 6,
				}, {
					Name:         "non-semver",
					Version:      "4fd2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 4,
				}, {
					Name:         "non-semver",
					Version:      "6343.4r32.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 3,
				}, {
					Name:         "non-semver",
					Version:      "5fd2.1212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortPackages(tt.pkgs)
			for i := range got {
				if diff := cmp.Diff(got[i].Version, tt.want[i].Version); diff != "" {
					t.Errorf("TestSortPackages() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
