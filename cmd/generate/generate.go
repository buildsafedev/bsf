package generate

import (
	"github.com/spf13/cobra"
)

// GenCmd represents the generate command
var GenCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate resolves dependencies and generates the Nix files",
	Long: `generate resolves dependencies and generates the Nix files. It is similiar to "go mod tidy"
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// sc, err := search.NewClientWithURL("https://api.history.nix-packages.com/")
		// if err != nil {
		// 	os.Exit(1)
		// }

		// m := model{sc: sc}
		// m.resetSpinner()
		// if _, err := tea.NewProgram(m).Run(); err != nil {
		// 	os.Exit(1)
		// }
	},
}
