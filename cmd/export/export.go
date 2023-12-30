package export

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/spf13/cobra"
)

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports Nix package outputs to an artifact(ex- OCI image)",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
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

		envNames := make([]string, 0, len(conf.Export))
		for _, ec := range conf.Export {
			envNames = append(envNames, ec.Environment)
			errStr := ec.Validate()
			if errStr != nil {
				fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Config for export block %s is invalid\n Error: %s", ec.Name, *errStr)))
				os.Exit(1)
			}
		}

		// todo: we could provide a TUI list dropdown to select
		if len(args) < 1 {
			fmt.Println(styles.ErrorStyle.Render("error: ", "export name is required. Valid export environments are: ", strings.Join(envNames, ", ")))
			os.Exit(1)
		}

		// check if env exists
		var found bool
		for _, env := range envNames {
			if env == args[0] {
				found = true
				break
			}
		}
		if !found {
			fmt.Println(styles.ErrorStyle.Render("error: ", "No such Environment found. Valid export environments are: ", strings.Join(envNames, ", ")))
		}

	},
}

func init() {
	// ExportCmd.Flags().StringVarP()
}
