package precheck

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/build"
	"github.com/buildsafedev/bsf/pkg/nix/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

const (
	nixversion = `v2.18.1`
)

// PreCheckCmd represents the precheck command that checks the pre-requisites
var PreCheckCmd = &cobra.Command{
	Use:   "precheck",
	Short: "precheck checks the pre-requisites for the bsf ",
	Long:  `precheck checks the current nix version and is flakes enabled? and various other pre-requisites for the bsf to work properly.`,
	Run: func(cmd *cobra.Command, args []string) {
		AllPrechecks()
	},
}

func checkVersionGreater(currentVer string, nixversion string) bool {
	val := semver.Compare(currentVer, nixversion)

	if val < 0 {
		return false
	}

	return true
}

// ValidateNixVersion checks if the current nix version is it compatible with bsf
func ValidateNixVersion() {
	currentver, err := cmd.NixVersion()
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error fetching nix version:", err.Error()))
		os.Exit(1)
	}

	if !checkVersionGreater(currentver, nixversion) {
		fmt.Println(styles.ErrorStyle.Render("Nix version should be", nixversion, "or above"))
		os.Exit(1)
	}

	fmt.Println(styles.HelpStyle.Render(" ✅ Nix version is", currentver))
}

// IsFlakesEnabled checks if the flakes are enabled in the nix configuration
func IsFlakesEnabled() {
	config, err := cmd.NixShowConfig()
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error fetching nix config:", err.Error()))
		os.Exit(1)
	}

	expectedKey := "experimental-features"

	value := config[expectedKey]
	if strings.Contains(value, "flakes") {
		fmt.Println(styles.HelpStyle.Render(" ✅ Flakes are enabled"))
	} else {
		fmt.Println(styles.ErrorStyle.Render("Flakes are not enabled"))
		os.Exit(1)
	}
}

func IsContainerDStoreEnabled() {
	conf, err := build.ReadDockerDaemonCfg()
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("err:", err.Error()))
		os.Exit(1)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(conf.Features), &data); err != nil {
		fmt.Println(styles.ErrorStyle.Render("err:", err.Error()))
		os.Exit(1)
	}

	features, _ := data["features"]

	if containerdSnapshotter, ok := features.(map[string]interface{})["containerd-snapshotter"]; ok {
		if value, ok := containerdSnapshotter.(bool); ok && value {
			fmt.Println(styles.HelpStyle.Render(" ✅ containerd-snapshotter is set to true"))
			return
		}
	}
	fmt.Println(styles.HelpStyle.Render(" ⚠️  containerd image store is not enabled [ https://docs.docker.com/storage/containerd/ ]"))

}

// AllPrechecks runs all the prechecks
func AllPrechecks() {
	fmt.Println(styles.TextStyle.Render("Running prechecks..."))
	ValidateNixVersion()
	IsFlakesEnabled()
	IsContainerDStoreEnabled()
	fmt.Println(styles.SucessStyle.Render(" Prechecks ran successfully"))
}
