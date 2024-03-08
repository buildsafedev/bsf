package search

import (
	"testing"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/vulnerability"
)

func TestDeriveAV(t *testing.T) {
	tests := []struct {
		name   string
		vector string
		want   string
	}{
		{
			name:   "Test Case 1",
			vector: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Network",
		},
		{
			name:   "Test Case 2",
			vector: "CVSS:3.1/AV:A/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Adjacent Network",
		},

		{
			name:   "Test Case 3",
			vector: "CVSS:3.1/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "",
		},
		{
			name:   "Test Case 4",
			vector: "CVSS:3.1/AV:L/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Local",
		},
		{
			name:   "Test Case 5",
			vector: "CVSS:3.1/AV:P/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   "Physical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vulnerability.DeriveAV(tt.vector); got != tt.want {
				t.Errorf("deriveAV() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := search.SortPackages(tt.pkgs); got[0].Version != tt.want[0].Version {
				t.Errorf("SortPackages() = %v, want %v", got, tt.want)
			}
		})
	}
}
