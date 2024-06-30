package template

import (
	"os"
	"testing"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

func TestTemplateMainForGolang(t *testing.T) {

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
		GoModule: &hcl2nix.GoModule{
			Name: "go-project",
		},
	}

	err := GenerateFlake(flake, os.Stdout, conf, false, "")
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

	err := GenerateFlake(flake, os.Stdout, conf, false, "")
	if err != nil {
		t.Error()
		t.FailNow()
	}
}

func TestTemplateMainForPython(t *testing.T) {

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
		PoetryApp: &hcl2nix.PoetryApp{
			ProjectDir: "../.",
		},
	}

	err := GenerateFlake(flake, os.Stdout, conf, false, "")
	if err != nil {
		t.Error()
		t.FailNow()
	}
}

func TestTemplateMainForJavacript(t *testing.T) {

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
		JsNpmApp: &hcl2nix.JsNpmApp{
			PackageName: "npm-project",
			PackageRoot: "./.",
		},
	}

	err := GenerateFlake(flake, os.Stdout, conf, false, "")
	if err != nil {
		t.Error()
		t.FailNow()
	}
}
