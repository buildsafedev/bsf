package develop

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	binit "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// DevCmd represents the Develop command
var DevCmd = &cobra.Command{
	Use:   "develop",
	Short: "develop spawns a development shell",
	Long: `develop spawns a development shell. All packages mentioned in bsf.hcl in development attribute will be available in the shell.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		sc, fh, err := binit.GetBSFInitializers()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = generate.Generate(fh, sc, nil)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = bgit.Add("bsf/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = nixcmd.Develop()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
	},
}
