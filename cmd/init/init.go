package init

import (
	"io/fs"
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

		sc, err := search.NewClientWithURL("https://api.history.nix-packages.com/")
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

func createBsfDirectory() ([]fs.DirEntry, error) {
	// check if the directory exists
	files, err := os.ReadDir("bsf")
	if err != nil {
		// check if the error is because the directory doesn't exist
		if os.IsNotExist(err) {
			err = os.Mkdir("bsf", 0755)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}
