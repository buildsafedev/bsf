package develop

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// DevCmd represents the Develop command
var DevCmd = &cobra.Command{
	Use:   "develop",
	Short: "develop spawns a development shell",
	Long: `develop spawns a development shell. All packages mentioned in bsf.hcl in development attribute will be available in the shell.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := nixcmd.Develop()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
	},
}
