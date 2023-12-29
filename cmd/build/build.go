package build

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Long: `builds the project based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut, err := nixcmd.Build()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: %v", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render(stdOut))

	},
}
