/*
Copyright © 2023 BuildSafe
*/

// Package cmd contains the main command for the bsf cli
package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	initCmd "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/search"
	"github.com/elewis787/boa"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bsf",
	Short: "bsf ",
	Long:  ` `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetHelpFunc(boa.HelpFunc)
	rootCmd.SetUsageFunc(boa.UsageFunc)

	rootCmd.AddCommand(initCmd.InitCmd)
	rootCmd.AddCommand(search.SearchCmd)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}

}
