package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Long: `builds the project based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files in version control system(ex: Git) before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("bsf.hcl"); err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: Has the project been initialized?"))
			fmt.Println(styles.HintStyle.Render("hint: ", "run `bsf init` to initialize the project"))
			os.Exit(1)
		}
		fmt.Println(styles.HighlightStyle.Render("Building, please be patient..."))

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		err = nixcmd.Build(conf)
		if err != nil {
			if isNoFileError(err.Error()) {
				fmt.Println(styles.ErrorStyle.Render(err.Error() + "\n Please ensure all necessary files are added/committed in your version control system"))
				fmt.Println(styles.HintStyle.Render("hint: run git add .  "))
				os.Exit(1)
			}

			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))

			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render("Build completed successfully, please check the result directory"))

	},
}

func isNoFileError(err string) bool {
	return strings.Contains(err, "No such file or directory")
}
