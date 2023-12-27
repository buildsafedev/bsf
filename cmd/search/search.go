package search

import (
	"fmt"
	"os"
	"sort"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// SearchCmd represents the init command
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search searches for packages",
	Long: `Search for Nix packages on the Nixpkgs repository
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", "package name is required").Error()))
			os.Exit(1)
		}

		sc, err := search.NewClientWithURL("https://api.history.nix-packages.com/")
		if err != nil {
			os.Exit(1)
		}

		packages, err := sc.ListPackageVersions(args[0])
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
			os.Exit(1)
		}

		m := model{packageOptionModel: packageOptionModel{},
			versionList: list.New(convertPackagesToItems(sortPackagesWithTimestamp(packages)),
				list.NewDefaultDelegate(), 0, 0)}
		m.versionList.Title = args[0]
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			os.Exit(1)
		}
	},
}

func convertPackagesToItems(packages []search.Package) []list.Item {
	items := make([]list.Item, 0, len(packages))
	for _, pkg := range packages {
		items = append(items, item{pkg: pkg})
	}

	return items
}

func sortPackagesWithTimestamp(packages []search.Package) []search.Package {
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].DateTime.After(packages[j].DateTime)
	})
	return packages
}
