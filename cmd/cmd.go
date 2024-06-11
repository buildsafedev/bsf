/*
Copyright © 2023 BuildSafe
*/

// Package cmd contains the main command for the bsf cli
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/elewis787/boa"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/attestation"
	"github.com/buildsafedev/bsf/cmd/build"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/develop"
	"github.com/buildsafedev/bsf/cmd/direnv"
	"github.com/buildsafedev/bsf/cmd/dockerfile"
	initCmd "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/nixgenerate"
	"github.com/buildsafedev/bsf/cmd/oci"
	"github.com/buildsafedev/bsf/cmd/precheck"
	"github.com/buildsafedev/bsf/cmd/scan"
	"github.com/buildsafedev/bsf/cmd/search"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/cmd/update"
)

var (
	// DebugDir is the directory where bsf project needs to be debugged
	DebugDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bsf",
	Short: "bsf CLI lets you manage OS dependencies of your application seamlessly",
	Long:  `Opinionated app dependency management tool.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	debugDir := getDebugPath()
	if debugDir != "" {
		err := os.Chdir(debugDir)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
	}

	rootCmd.SetHelpFunc(boa.HelpFunc)
	rootCmd.SetUsageFunc(boa.UsageFunc)

	rootCmd.AddCommand(initCmd.InitCmd)
	rootCmd.AddCommand(search.SearchCmd)
	rootCmd.AddCommand(develop.DevCmd)
	rootCmd.AddCommand(build.BuildCmd)
	rootCmd.AddCommand(scan.ScanCmd)
	rootCmd.AddCommand(update.UpdateCmd)
	rootCmd.AddCommand(attestation.AttCmd)
	rootCmd.AddCommand(direnv.Direnv)

	if os.Getenv("BSF_DEBUG_MODE") == "true" {
		rootCmd.AddCommand(configure.ConfigureCmd)
		rootCmd.AddCommand(nixgenerate.NixGenCmd)
		rootCmd.AddCommand(precheck.PreCheckCmd)
	}
	rootCmd.AddCommand(oci.OCICmd)
	rootCmd.AddCommand(dockerfile.DFCmd)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}

}

func getDebugPath() string {
	if os.Getenv("BSF_DEBUG_DIR") != "" {
		DebugDir = os.Getenv("BSF_DEBUG_DIR")
	}
	return DebugDir
}
