package dockerfile

import (
	"fmt"
	"os"

	bsfinit "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/df"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init a Dockerfile into your project",
	Run: func(cmd *cobra.Command, args []string) {

		dfType, err := bsfinit.IoPrompt("Which language is your app in?(Go/Python/Rust)")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		isHermetic, err := bsfinit.YesNoPrompt("Would you like a hermetic build? This requires you to vendor your dependencies(Recommended: yes)")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		dfFile, err := os.Create("Dockerfile")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		defer dfFile.Close()

		err = df.GenerateDF(dfFile, dfType, isHermetic)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
	},
}
