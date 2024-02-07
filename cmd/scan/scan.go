package scan

import (
	"context"
	"fmt"
	"os"
	"strings"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
)

// ScanCmd represents the build command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans the given package name and version for vulnerabilities.",
	Long: `scans the given package name and version for vulnerabilities.
	Example : bsf scan <package name:package version> / <package name> <package version>
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			fmt.Println(styles.ErrorStyle.Render(fmt.Errorf("error: %v, received %v", "accepts between 1 and 2 arg(s)", len(args)).Error()))
			os.Exit(1)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var name, version string

		if len(args) == 1 {
			nameWithVersion := strings.SplitN(args[0], ":", 2)

			if len(nameWithVersion) < 2 {
				fmt.Println(styles.ErrorStyle.Render(fmt.Errorf("error: %v", "version is required").Error()))
				os.Exit(1)
			}

			name = nameWithVersion[0]
			version = nameWithVersion[1]
		}

		if len(args) == 2 {
			name = args[0]
			version = args[1]
		}

		fmt.Println(styles.BaseStyle.Render("info: ", "Scanning..."))

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			os.Exit(1)
		}

		vulnerabilities, err := sc.FetchVulnerabilities(context.Background(), &bsfv1.FetchVulnerabilitiesRequest{
			Name:    name,
			Version: version,
		})
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		m := initVulnTable(vulnerabilities)
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(errorStyle.Render(fmt.Errorf("error: %v", err).Error()))
			os.Exit(1)
		}
	},
}
