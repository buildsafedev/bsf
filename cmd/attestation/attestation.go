package attestation

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/spf13/cobra"
)

// AttCmd represents the attestation command
var AttCmd = &cobra.Command{
	Use:   "att",
	Short: "perform attestation ops",
	Long:  `used to perform various operations on your attestations`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styles.HintStyle.Render("hint: use bsf att with a subcomand"))
		os.Exit(1)
	},
}

func init() {
	// add subcommand to list predicates
	AttCmd.AddCommand(listCmd)
	// add subcommand to print predicates
	AttCmd.AddCommand(catCmd)
}
