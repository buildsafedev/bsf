package export

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/build"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

var (
	output string
)

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "exports Nix package outputs to an artifact(ex- OCI image)",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// todo: we could provide a TUI list dropdown to select
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

		envNames := make([]string, 0, len(conf.Export))
		var found bool
		env := hcl2nix.ExportConfig{}
		for _, ec := range conf.Export {
			errStr := ec.Validate()
			if errStr != nil {
				fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Config for export block %s is invalid\n Error: %s", ec.Name, *errStr)))
				os.Exit(1)
			}

			if ec.Environment == args[0] {
				found = true
				env = ec
				break
			}
			envNames = append(envNames, ec.Environment)
		}

		if !found {
			fmt.Println(styles.ErrorStyle.Render("error:", "No such environment found. Valid export environments are:", strings.Join(envNames, ", ")))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf export <environment name>` to export the environment"))
			os.Exit(1)
		}

		if output != "" {
			fh, err := os.Create(output)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			defer fh.Close()

			err = build.GenerateDockerfile(fh, env)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			return
		}
		err = build.Build(env)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

	},
}

func init() {
	ExportCmd.Flags().StringVarP(&output, "output", "o", "", "location of the generated Dockerfile")
}
