package hcl2nix

import (
	"bytes"
	"io"
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
