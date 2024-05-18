package attestation

import (
	"github.com/spf13/cobra"
)

// AttCmd represents the attestation command
var AttCmd = &cobra.Command{
	Use:   "att",
	Short: "perform attestation ops",
	Long: `used to perform various operations on your attestations`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	// add subcommand to list predicates
	AttCmd.AddCommand(listCmd)
}
