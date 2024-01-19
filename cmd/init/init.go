package init

import (
	"os"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "init setups package management for the project",
	Long: `init setups package management for the project. It setups Nix files based on the language detected.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		sc, err := search.NewClientWithAddr("localhost:8080")
		if err != nil {
			os.Exit(1)
		}

		m := model{sc: sc}
		m.resetSpinner()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			os.Exit(1)
		}
	},
}
