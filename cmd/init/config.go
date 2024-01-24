package init

import (
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

func generatehcl2NixConf(pt langdetect.ProjectType, pd *langdetect.ProjectDetails) hcl2nix.Config {
	switch pt {
	case langdetect.GoModule:
		var name, entrypoint string
		if pd != nil {
			name = pd.Name
			entrypoint = pd.Entrypoint
			if entrypoint == "" {
				entrypoint = "./."
			}

		}
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{
				Development: []string{"go@1.21.4", "gotools@0.7.0", "delve@1.21.2"},
				// todo: maybe we should dynamically inject the latest version of such runtime packages(cacert)?
				Runtime: []string{"cacert@3.95"},
			},
			GoModule: &hcl2nix.GoModule{
				Name:       name,
				SourcePath: entrypoint,
			},
		}
	default:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{},
		}
	}
}
