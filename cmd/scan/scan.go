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
	Example : bsf scan <package name> <package version>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.ErrorStyle.Render(fmt.Errorf("error: %v", "package name and version is required").Error()))
			os.Exit(1)
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

		nameWithVersion := strings.Split(args[0], ":")

		vulnerabilities, err := sc.FetchVulnerabilities(context.Background(), &bsfv1.FetchVulnerabilitiesRequest{
			Name:    nameWithVersion[0],
			Version: nameWithVersion[1],
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
