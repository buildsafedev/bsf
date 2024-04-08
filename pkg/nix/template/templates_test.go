package template

import (
	"os"
	"testing"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

func TestTemplateMain(t *testing.T) {

	flake := Flake{
		Description: "Simple flake",
		NixPackageRevisions: []string{
			"a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"a9bf124c46ef298113270b1f84a164865987a91c",
		},
		DevPackages: map[string]string{
			"gotools": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"go_1_19": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
		},
		RuntimePackages: map[string]string{
			"bash": "a9bf124c46ef298113270b1f84a164865987a91c",
		},
	}

	err := GenerateFlake(flake, os.Stdout, nil)
	if err != nil {
		t.Error()
		t.FailNow()
	}
}

func TestTemplateMainForRust(t *testing.T) {

	flake := Flake{
		Description: "Simple flake",
		NixPackageRevisions: []string{
			"a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"a9bf124c46ef298113270b1f84a164865987a91c",
		},
		DevPackages: map[string]string{
			"gotools": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"go_1_19": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
		},
		RuntimePackages: map[string]string{
			"bash": "a9bf124c46ef298113270b1f84a164865987a91c",
		},
	}

	conf := &hcl2nix.Config{
		RustApp: &hcl2nix.RustApp{
			RustVersion: "1.25.0",
			Release:     true,
		},
	}

	err := GenerateFlake(flake, os.Stdout, conf)
	if err != nil {
		t.Error()
		t.FailNow()
	}
}
