package hcl2nix

import (
	"bytes"
	"io"
	"slices"
	"testing"

	bstrings "github.com/buildsafedev/bsf/pkg/strings"
	"github.com/buildsafedev/bsf/pkg/update"
)

func TestWriteConfig(t *testing.T) {
	err := WriteConfig(Config{
		Packages: Packages{
			Development: []string{"go@~1.19.3", "nodejs"},
			Runtime:     []string{"python", "ruby"},
		},
	}, io.Discard)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

}

func TestReadConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteConfig(Config{
		Packages: Packages{
			Development: []string{"go@~1.19.3", "nodejs"},
			Runtime:     []string{"python", "ruby"},
		},
	}, buf)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	config, err := ReadConfig(buf.Bytes(), io.Discard)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(config.Packages.Development) != 2 || len(config.Packages.Runtime) != 2 {
		t.Error("Packages not read correctly")
		t.Fail()
	}

}

func TestPreferNewElemenets(t *testing.T) {
	tests := []struct {
		name           string
		existingConfig Config
		newPackages    Config
		result         Config
	}{
		{
			name: "development duplicates",
			existingConfig: Config{
				Packages: Packages{
					Development: []string{"go@~1.19.3", "nodejs@1.1"},
				},
			},
			newPackages: Config{
				Packages: Packages{
					Development: []string{"go@~1.20.3"},
				},
			},
			result: Config{
				Packages: Packages{
					Development: []string{"go@~1.20.3", "nodejs@1.1"},
				},
			},
		},
		{
			name: "runtime duplicates",
			existingConfig: Config{
				Packages: Packages{
					Runtime: []string{"go@~1.20.11", "ruby"},
				},
			},
			newPackages: Config{
				Packages: Packages{
					Runtime: []string{"go@1.20.11"},
				},
			},
			result: Config{
				Packages: Packages{
					Runtime: []string{"go@1.20.11", "ruby"},
				},
			},
		},

		{
			name: "both duplicates",
			existingConfig: Config{
				Packages: Packages{
					Development: []string{"go@~1.20.11", "nodejs"},
					Runtime:     []string{"python", "go@~1.18.11"},
				},
			},
			newPackages: Config{
				Packages: Packages{
					Development: []string{"go@~1.19.3"},
					Runtime:     []string{"go@~1.19.3"},
				},
			},
			result: Config{
				Packages: Packages{
					Development: []string{"go@~1.19.3", "nodejs"},
					Runtime:     []string{"python", "go@~1.19.3"},
				},
			},
		},
		{
			name: "add new package (no duplicates)",
			existingConfig: Config{
				Packages: Packages{
					Development: []string{"go@~1.20.11", "nodejs"},
					Runtime:     []string{"python", "go@~1.18.11"},
				},
			},
			newPackages: Config{
				Packages: Packages{
					Development: []string{"go-task@~1.3.3"},
					Runtime:     []string{"go-task@~1.3.3"},
				},
			},
			result: Config{
				Packages: Packages{
					Development: []string{"go@~1.20.11", "nodejs", "go-task@~1.3.3"},
					Runtime:     []string{"python", "go@~1.18.11", "go-task@~1.3.3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parse := func(s string) string {
				name, _ := update.ParsePackage(s)
				return name
			}

			dev := bstrings.PreferNewSliceElements(tt.existingConfig.Packages.Development, tt.newPackages.Packages.Development, parse)
			runtime := bstrings.PreferNewSliceElements(tt.existingConfig.Packages.Runtime, tt.newPackages.Packages.Runtime, parse)

			for _, d := range tt.result.Packages.Development {
				if !slices.Contains(dev, d) {
					t.Errorf("Pkg %s not  found in development env", d)
					t.Fail()
				}
			}

			for _, r := range tt.result.Packages.Runtime {
				if !slices.Contains(runtime, r) {
					t.Errorf("Pkg %s not found in runtime env", r)
					t.Fail()
				}
			}
		})
	}
}
