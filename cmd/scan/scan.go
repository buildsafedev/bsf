package scan

import (
	"context"
	"fmt"
	"os"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/spf13/cobra"
)

// ScanCmd represents the build command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans the given package name and version for vulnerabilities.",
	Long: `scans the given package name and version for vulnerabilities.
	Example : bsf scan <package name> <package version>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
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

		vulnerabilies, err := sc.FetchVulnerabilities(context.Background(), &bsfv1.FetchVulnerabilitiesRequest{
			Name:    args[0],
			Version: args[1],
		})
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("%d vulnerabilities found", len(vulnerabilies.GetVulnerabilities()))))
		// TODO: table view of vulnerabilities

	},
}
