package search

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/pkg/clients/search"
	buildsafev1 "github.com/buildsafedev/cloud-api/apis/v1"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

var sc buildsafev1.SearchServiceClient

// SearchCmd represents the init command
var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search searches for packages",
	Long: `Search for Nix packages on the Nixpkgs repository
	`,
	Example: `bsf search <package name>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			// todo: now that we have server side pagination, we can support ad-hoc browsing of packages
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", "package name is required").Error()))
			os.Exit(1)
		}
		if os.Getenv("BSF_DEBUG") != "" && strings.ToLower(os.Getenv("BSF_DEBUG")) == "true" {
			// todo: maybe we should have an option to send this to server to help debug?
			if f, err := tea.LogToFile("debug.log", "help"); err != nil {
				fmt.Println("Couldn't open a file for logging:", err)
				os.Exit(1)
			} else {
				defer func() {
					err = f.Close()
					if err != nil {
						log.Fatal(err)
					}
				}()
			}
		}

		var err error
		sc, err = search.NewClientWithAddr("localhost:8080")
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", err).Error()))
			os.Exit(1)
		}

		packages, err := sc.ListPackages(cmd.Context(), &buildsafev1.ListPackagesRequest{
			// since we have a term specified here, assuming that the result won't larger than 1k items
			PageSize:   999,
			PageToken:  1,
			SearchTerm: args[0],
		})
		if err != nil {
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", err).Error()))
			os.Exit(1)
		}
		items := convLPR2Items(packages)
		m := InitSearch(items)
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", err).Error()))
			os.Exit(1)
		}

	},
}
