package dockerfile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
)

var (
	platform string
)

func init() {
	DFCmd.Flags().StringVarP(&platform, "platform", "p", "", "The platform to build the image for")
	DFCmd.AddCommand(DGCmd)
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
		fmt.Println(styles.HintStyle.Render("hint: use bsf dockerfile with a subcomand"))
		os.Exit(1)
	},
}
