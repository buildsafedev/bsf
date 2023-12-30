/*
Copyright © 2023 BuildSafe
*/

// Package cmd contains the main command for the bsf cli
package cmd

import (
	"context"
	"os"

	"github.com/elewis787/boa"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/build"
	"github.com/buildsafedev/bsf/cmd/develop"
	"github.com/buildsafedev/bsf/cmd/export"
	"github.com/buildsafedev/bsf/cmd/generate"
	initCmd "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/run"
	"github.com/buildsafedev/bsf/cmd/search"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bsf",
	Short: "bsf CLI ",
	Long:  `An opinionated application definition framework`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetHelpFunc(boa.HelpFunc)
	rootCmd.SetUsageFunc(boa.UsageFunc)

	rootCmd.AddCommand(initCmd.InitCmd)
	rootCmd.AddCommand(search.SearchCmd)
	rootCmd.AddCommand(generate.GenCmd)
	rootCmd.AddCommand(develop.DevCmd)
	rootCmd.AddCommand(build.BuildCmd)
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(export.ExportCmd)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}

}
