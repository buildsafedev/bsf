package update

import (
	"testing"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

func TestGetLatestPatchVersion(t *testing.T) {
	tests := []struct {
		name     string
		response *buildsafev1.FetchPackagesResponse
		version  string
		want     string
	}{
		{
			name: "Test Case 1",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.2.1"},
					{Version: "1.2.2"},
					{Version: "1.2.3"},
				},
			},
			version: "1.2.0",
			want:    "1.2.3",
		},
		{
			name: "Test Case 2",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.3.1"},
					{Version: "1.3.2"},
					{Version: "1.3.3"},
				},
			},
			version: "1.3.0",
			want:    "1.3.3",
		},
		{
			name:     "Test Case 3",
			response: nil,
			version:  "1.3.0",
			want:     "",
		},
		{
			name: "Test Case 4",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.3.1"},
					{Version: "1.3.2"},
					{Version: "1.3.3"},
					{Version: "1.4.0"},
				},
			},
			version: "1.4.0",
			want:    "1.4.0",
		},

		{
			name: "Test Case 5",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.3.1"},
					{Version: "1.3.2"},
					{Version: "1.3.3"},
				},
			},
			version: "1.3.3",
			want:    "1.3.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLatestPatchVersion(tt.response, tt.version); got != tt.want {
				t.Errorf("getLatestPatchVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLatestMinorVersion(t *testing.T) {
	tests := []struct {
		name     string
		response *buildsafev1.FetchPackagesResponse
		version  string
		want     string
	}{
		{
			name: "Test Case 1",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.2.1"},
					{Version: "1.3.2"},
					{Version: "1.4.3"},
				},
			},
			version: "1.2.0",
			want:    "1.4.3",
		},
		{
			name: "Test Case 2",
			response: &buildsafev1.FetchPackagesResponse{
				Packages: []*buildsafev1.Package{
					{Version: "1.3.0"},
					{Version: "1.4.1"},

					{Version: "2.3.2"},
					{Version: "3.3.3"},
				},
			},
			version: "1.3.0",
			want:    "1.4.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLatestMinorVersion(tt.response, tt.version); got != tt.want {
				t.Errorf("getLatestPatchVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePackage(t *testing.T) {
	tests := []struct {
		name     string
		pkg      string
		wantName string
		wantVer  string
	}{
		{
			name:     "Test Case 1",
			pkg:      "go@^1.20",
			wantName: "go",
			wantVer:  "1.20",
		},
		{
			name:     "Test Case 2",
			pkg:      "go@~1.20",
			wantName: "go",
			wantVer:  "1.20",
		},
		{
			name:     "Test Case 3",
			pkg:      "go@1.20",
			wantName: "go",
			wantVer:  "1.20",
		},
		{
			name:     "Test Case 4",
			pkg:      "go",
			wantName: "",
			wantVer:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotVer := ParsePackage(tt.pkg)
			if gotName != tt.wantName || gotVer != tt.wantVer {
				t.Errorf("ParsePackage() = %v, %v, want %v, %v", gotName, gotVer, tt.wantName, tt.wantVer)
			}
		})
	}
}

func TestParseUpdateType(t *testing.T) {
	tests := []struct {
		name string
		pkg  string
		want int
	}{
		{
			name: "Test Case 1",
			pkg:  "go@^1.20",
			want: UpdateTypeMinor,
		},
		{
			name: "Test Case 2",
			pkg:  "go@~1.20",
			want: UpdateTypePatch,
		},
		{
			name: "Test Case 3",
			pkg:  "go@1.20",
			want: UpdateTypePinned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseUpdateType(tt.pkg); got != tt.want {
				t.Errorf("ParseUpdateType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "Test Case 1 - Equal slices",
			a:    []string{"pkg1@~v1.0.0", "pkg2@~v1.1.0", "pkg3@~v1.2.0"},
			b:    []string{"pkg3@~v1.2.0", "pkg2@~v1.1.0", "pkg1@~v1.0.0"},
			want: true,
		},
		{
			name: "Test Case 2 - Different slices",
			a:    []string{"pkg1@~v1.0.0", "pkg2@~v1.1.0", "pkg3@~v1.2.0"},
			b:    []string{"pkg1@~v1.2.0", "pkg2@~v1.1.0", "pkg3@~v1.3.0"},
			want: false,
		},
		{
			name: "Test Case 3 - Different lengths",
			a:    []string{"pkg1@~v1.0.0", "pkg2@~v1.1.0"},
			b:    []string{"pkg2@~v1.2.0", "pkg1@~v1.1.0", "pkg3@~v1.0.0"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareVersions(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
