package hcl2nix

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestWriteConfig(t *testing.T) {
	err := WriteConfig(Config{
		Packages: Packages{
			Development: []string{"go", "nodejs"},
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
			Development: []string{"go", "nodejs"},
			Runtime:     []string{"python", "ruby"},
		},
	}, buf)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	config, err := ReadConfig(buf.Bytes())
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(config.Packages.Development) != 2 || len(config.Packages.Runtime) != 2 {
		t.Error("Packages not read correctly")
		t.Fail()
	}

}

func TestAddStringToSet(t *testing.T) {
	tests := []struct {
		name             string
		existingPackages []string
		newPackages      []string
		expected         []string
	}{
		{
			name:             "no duplicates",
			existingPackages: []string{"package1", "package2"},
			newPackages:      []string{"package3", "package4"},
			expected:         []string{"package1", "package2", "package3", "package4"},
		},
		{
			name:             "with duplicates",
			existingPackages: []string{"package1", "package2"},
			newPackages:      []string{"package2", "package3"},
			expected:         []string{"package1", "package2", "package3"},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.existingPackages = addStringToSet(tt.existingPackages, tt.newPackages)
			if !reflect.DeepEqual(tt.existingPackages, tt.expected) {
				t.Errorf("got %v, want %v", tt.existingPackages, tt.expected)
			}
		})
	}
}
