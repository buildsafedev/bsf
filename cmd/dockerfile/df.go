package dockerfile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
)

func init() {
	DFCmd.AddCommand(DGCmd)
}

// DFCmd represents the generate command
var DFCmd = &cobra.Command{
	Use:     "dockerfile",
	Short:   "dockerfile generates a dockerfile for the app",
	Aliases: []string{"df"},
	Long: `
	bsf dockerfile <artifact> 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styles.HintStyle.Render("hint: use bsf dockerfile with a subcomand"))
		os.Exit(1)
	},
}
