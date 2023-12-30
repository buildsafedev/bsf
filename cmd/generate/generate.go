package generate

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// GenCmd represents the generate command
var GenCmd = &cobra.Command{
	Use:     "generate",
	Short:   "generate resolves dependencies and generates the Nix files",
	Aliases: []string{"gen"},
	Long: `generate resolves dependencies and generates the Nix files. It is similiar to "go mod tidy"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("bsf.hcl"); err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()+"\nHas the project been initialized?"))
			fmt.Println(styles.HintStyle.Render("hint: ", "run `bsf init` to initialize the project"))
			os.Exit(1)
		}

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
