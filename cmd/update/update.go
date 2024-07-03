package update

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/generate"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/update"
)

var updateCmdOptions struct {
	check bool
}

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates dependencies to highest version based on constraints",
	Long: `Updates can be done for development and runtime dependencies based on constraints. Following constraints are supported:
		~ : latest patch version
		^ : latest minor version

		Currently, only packages following semver versioning are supported.

		`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styles.TextStyle.Render("Updating..."))

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		data, err := os.ReadFile("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		var dstErr bytes.Buffer
		hconf, err := hcl2nix.ReadConfig(data, &dstErr)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(dstErr.String()))
			os.Exit(1)
		}

		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		devVersionMap := fetchPackageVersions(hconf.Packages.Development, sc)
		runtimeVersionMap := fetchPackageVersions(hconf.Packages.Runtime, sc)

		devVersions := parsePackagesForUpdates(devVersionMap)
		runtimeVersions := parsePackagesForUpdates(runtimeVersionMap)

		newPackages := hcl2nix.Packages{
			Development: devVersions,
			Runtime:     runtimeVersions,
		}

		if updateCmdOptions.check {
			if !compareVersions(hconf.Packages.Development, devVersions) || !compareVersions(hconf.Packages.Runtime, runtimeVersions) {
				fmt.Println(styles.WarnStyle.Render("Updates are available"))
				os.Exit(1)
			} else {
				fmt.Println(styles.SucessStyle.Render("No updates available"))
				os.Exit(0)
			}
		}

		fh, err := hcl2nix.NewFileHandlers(true)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error creating file handlers: %s", err.Error()))
			os.Exit(1)
		}
		// changing file handler to allow writes
		fh.ModFile, err = os.Create("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error creating bsf.hcl: %s", err.Error()))
			os.Exit(1)
		}

		err = hcl2nix.SetPackages(data, newPackages, fh.ModFile)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error updating bsf.hcl: %s", err.Error()))
			os.Exit(1)
		}

		err = generate.Generate(fh, sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("Error generating files: %s", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render("Updated ran successfully"))
	},
}

// This function compares the devVersions and the runtimeVersions
func compareVersions(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	counts := make(map[string]int)
	for _, item := range a {
		counts[item]++
	}
	for _, item := range b {
		if counts[item] == 0 {
			return false
		}
		counts[item]--
	}

	for _, count := range counts {
		if count != 0 {
			return false
		}
	}

	return true
}

func parsePackagesForUpdates(versionMap map[string]*buildsafev1.FetchPackagesResponse) []string {
	newVersions := make([]string, 0, len(versionMap))

	for k, v := range versionMap {
		name, version := update.ParsePackage(k)

		if !semver.IsValid("v" + version) {
			fmt.Println(styles.WarnStyle.Render("warning:", "skipping ", k, " as only semver updates are supported currently"))
			newVersions = append(newVersions, k)
			continue
		}

		updateType := update.ParseUpdateType(k)

		switch updateType {
		case update.UpdateTypePatch:
			newVer := update.GetLatestPatchVersion(v, version)
			if newVer != "" {
				newVersions = append(newVersions, fmt.Sprintf("%s@~%s", name, newVer))
			}

		case update.UpdateTypeMinor:
			newVer := update.GetLatestMinorVersion(v, version)
			if newVer != "" {
				newVersions = append(newVersions, fmt.Sprintf("%s@^%s", name, newVer))
			}

		case update.UpdateTypePinned:
			newVersions = append(newVersions, k)
			continue
		}
	}
	return newVersions
}

func fetchPackageVersions(packages []string, sc buildsafev1.SearchServiceClient) map[string]*buildsafev1.FetchPackagesResponse {
	versionsMap := make(map[string]*buildsafev1.FetchPackagesResponse)

	var wg sync.WaitGroup
	for _, pkg := range packages {
		name, version := update.ParsePackage(pkg)
		if name == "" || version == "" {
			fmt.Println(styles.ErrorStyle.Render("error:", "invalid package name or version"))
			continue
		}

		wg.Add(1)
		go func(name, pkg string) {
			defer wg.Done()
			allVersions, err := sc.FetchPackages(context.Background(), &buildsafev1.FetchPackagesRequest{
				Name: name,
			})
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error finding ", name, ":", err.Error()))
				return
			}

			versionsMap[pkg] = allVersions
		}(name, pkg)
	}
	wg.Wait()

	return versionsMap
}

func init() {
	UpdateCmd.PersistentFlags().BoolVarP(&updateCmdOptions.check, "check", "c", false, "Check for updates without applying them")
}
