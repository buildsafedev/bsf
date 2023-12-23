package template

import "testing"

func TestTemplateMain(t *testing.T) {

	flake := Flake{
		Description: "Simple flake",
		PackageInputUrls: []string{
			"a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"a9bf124c46ef298113270b1f84a164865987a91c",
		},
		DevPackages: map[string]string{
			"gotools": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
			"go_1_19": "a89ba043dda559ebc57fc6f1fa8cf3a0b207f688",
		},
		PackageInputs: map[string]string{"go": "pkgs.go"},
		RuntimePackages: map[string]string{
			"bash": "a9bf124c46ef298113270b1f84a164865987a91c",
		},
	}

	_, err := templateMain(flake)
	if err != nil {
		t.Error()
		t.FailNow()
	}
}
