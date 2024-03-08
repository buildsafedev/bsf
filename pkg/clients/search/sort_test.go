package search

import (
	"testing"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

func TestSortPackages(t *testing.T) {
	tests := []struct {
		name string
		pkgs []*buildsafev1.Package
		want []*buildsafev1.Package
	}{
		{
			name: "Test Case 1",
			pkgs: []*buildsafev1.Package{
				{
					Name:         "test",
					Version:      "1.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
				{
					Name:         "test",
					Version:      "2.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
			},
			want: []*buildsafev1.Package{
				{
					Name:         "test",
					Version:      "2.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 2,
				},
				{
					Name:         "test",
					Version:      "1.0.0",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
			},
		},
		{
			name: "Test Case 2",
			pkgs: []*buildsafev1.Package{
				{
					Name:         "test",
					Version:      "234ca.b243.cc32c",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
				{
					Name:         "test",
					Version:      "3213.12212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
				{
					Name:         "test",
					Version:      "1.5.22",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
			},
			want: []*buildsafev1.Package{
				{
					Name:         "test",
					Version:      "234ca.b243.cc32c",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
				{
					Name:         "test",
					Version:      "3213.12212.1212",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
				{
					Name:         "test",
					Version:      "1.5.22",
					SpdxId:       "MIT",
					Free:         true,
					Homepage:     "https://test.com",
					EpochSeconds: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortPackages(tt.pkgs); got[0].Version != tt.want[0].Version {
				t.Errorf("SortPackages() = %v, want %v", got, tt.want)
			}
		})
	}
}
