package generate

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// GenCmd represents the generate command
var GenCmd = &cobra.Command{
	Use:     "generate",
	Short:   "generate resolves dependencies and generates the Nix files",
	Aliases: []string{"gen"},
	Long: `generate resolves dependencies and generates the Nix files.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf export <environment name>` to export the environment"))
			os.Exit(1)
		}

		data, err := os.ReadFile("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		var dstErr bytes.Buffer
		conf, err := hcl2nix.ReadConfig(data, &dstErr)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(dstErr.String()))
			os.Exit(1)
		}

		envNames := make([]string, 0, len(conf.OCIArtifact))
		var found bool
		// env := hcl2nix.OCIExportConfig{}
		// for _, ec := range conf.Export {
		// 	errStr := ec.Validate()
		// 	if errStr != nil {
		// 		fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Config for export block %s is invalid\n Error: %s", ec.Name, *errStr)))
		// 		os.Exit(1)
		// 	}

		// 	if ec.Environment == args[0] {
		// 		found = true
		// 		env = ec
		// 		break
		// 	}
		// 	envNames = append(envNames, ec.Environment)
		// }

		if !found {
			fmt.Println(styles.ErrorStyle.Render("error:", "No such environment found. Valid export environments are:", strings.Join(envNames, ", ")))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf export <environment name>` to export the environment"))
			os.Exit(1)
		}

	},
}
