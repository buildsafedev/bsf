package run

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// RunCmd represents the build command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "runs the project",
	Long: `runs the project after building it based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files in version control system(ex: Git) before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		err := nixcmd.Run()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: %v", err.Error()))
			os.Exit(1)
		}

	},
}
