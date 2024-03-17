package precheck

import (
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

// IsContainerDStoreEnabled checks is containerd storage is enabled
func IsContainerDStoreEnabled() {
	conf, err := build.GetSnapshotter()
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("err:", err.Error()))
		os.Exit(1)
	}
	resp := isSnapshotterEnabled(conf)
	if resp {
		fmt.Println(styles.HelpStyle.Render(" ✅ containerd-snapshotter is set to true"))
	} else {
		fmt.Println(styles.HelpStyle.Render(" ⚠️  Export functionality might not work as containerd image store is not enabled . To enable it, refer https://docs.docker.com/storage/containerd/ "))
	}
}

func isSnapshotterEnabled(conf string) bool {
	expectedOutput := " '[[driver-type io.containerd.snapshotter.v1]]' "
	if strings.Compare(strings.TrimSpace(expectedOutput), strings.TrimSpace(conf)) == 0 {
		return true
	} else {
		return false
	}
}

// AllPrechecks runs all the prechecks
func AllPrechecks() {
	fmt.Println(styles.TextStyle.Render("Running prechecks..."))
	ValidateNixVersion()
	IsFlakesEnabled()
	IsContainerDStoreEnabled()
	fmt.Println(styles.SucessStyle.Render(" Prechecks ran successfully"))
}
