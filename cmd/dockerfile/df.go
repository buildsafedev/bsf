package dockerfile

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	binit "github.com/buildsafedev/bsf/cmd/init"
	ocicmd "github.com/buildsafedev/bsf/cmd/oci"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/builddocker"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
)

var (
	output, platform string
)

func init() {
	DFCmd.Flags().StringVarP(&output, "output", "o", "", "location of the dockerfile generated")
	DFCmd.Flags().StringVarP(&platform, "platform", "p", "", "The platform to build the image for")
}

// DFCmd represents the generate command
var DFCmd = &cobra.Command{
	Use:     "dockerfile",
	Short:   "dockerfile generates a dockerfile for the app",
	Aliases: []string{"df"},
	Long: `
	bsf dockerfile <artifact> 
	bsf dockerfile <artifact> --platform <platform>
	bsf dockerfile <artifact> --platform <platform> --output <output filename>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf dockerfile <artifact>` to export the environment"))
			os.Exit(1)
		}

		env, p, err := ocicmd.ProcessPlatformAndConfig(platform, args[0])
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		platform = p

		sc, fh, err := binit.GetBSFInitializers()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = generate.Generate(fh, sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = bgit.Add("bsf/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		var dfw io.Writer
		if output == "" {
			dfw = os.Stdout
		} else {
			dfh, err := os.Create(output)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
				os.Exit(1)
			}
			defer dfh.Close()
			dfw = dfh
		}

		err = builddocker.GenerateDockerfile(dfw, env, platform)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

	},
}
