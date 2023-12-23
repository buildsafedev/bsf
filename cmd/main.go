/*
Copyright © 2023 BuildSafe
*/

// Package cmd contains the main command for the bsf cli
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	initCmd "github.com/buildsafedev/bsf/cmd/init"
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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd.InitCmd)
}
