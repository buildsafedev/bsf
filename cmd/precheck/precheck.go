package precheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/nix/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

const (
	nixversion = `v2.18.1`
)

var dockerdaemon_json = `/etc/docker/daemon.json`

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

// IsSnapshotterEnabled checks containerd image store is enabled. or not
func IsSnapshotterEnabled() (bool, error) {
	file, err := os.Open(dockerdaemon_json)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false, err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return false, err
	}

	if features, ok := data["features"].(map[string]interface{}); ok {
		if snapshotter, ok := features["containerd-snapshotter"]; ok {
			if snapshotterBool, ok := snapshotter.(bool); ok && snapshotterBool {
				fmt.Println(styles.HelpStyle.Render(" ✅ containerd-snapshotter is set to true"))
				return true, nil
			}
		}
	}

	return false, nil

}

// AllPrechecks runs all the prechecks
func AllPrechecks() {
	fmt.Println(styles.TextStyle.Render("Running prechecks..."))
	ValidateNixVersion()
	IsFlakesEnabled()
	resp, err := IsSnapshotterEnabled()
	if err != nil {
		fmt.Println(err)
	}
	if !resp {
		fmt.Println(styles.HelpStyle.Render(" ⚠️  containerd image store is not enabled [ https://docs.docker.com/storage/containerd/ ]"))
	}
	fmt.Println(styles.SucessStyle.Render(" Prechecks ran successfully"))
}
