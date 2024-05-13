package oci

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

var (
	platform string
)
var (
	supportedPlatforms = []string{"linux/amd64", "linux/arm64"}
)

// OCICmd represents the export command
var OCICmd = &cobra.Command{
	Use:   "oci",
	Short: "Builds an OCI image",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// todo: we could provide a TUI list dropdown to select
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf oci <environment name>` to build an OCI image"))
			os.Exit(1)
		}

		if platform == "" {
			tos := runtime.GOOS
			tarch := runtime.GOARCH
			platform = fmt.Sprintf("%s/%s", tos, tarch)
		}

		for _, sp := range supportedPlatforms {
			if strings.Contains(platform, sp) {
				break
			}
			fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Platform %s is not supported. Supported platforms are %s", platform, strings.Join(supportedPlatforms, ", "))))
			os.Exit(1)
		}

		data, err := os.ReadFile("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		var dstErr bytes.Buffer
		conf, err := hcl2nix.ReadConfig(data, &dstErr)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(dstErr.String()))
			os.Exit(1)
		}

		envNames := make([]string, 0, len(conf.OCIArtifact))
		var found bool
		env := hcl2nix.OCIArtifact{}
		for _, ec := range conf.OCIArtifact {
			errStr := ec.Validate(conf)
			if errStr != nil {
				fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Config for export block %s is invalid\n Error: %s", ec.Name, *errStr)))
				os.Exit(1)
			}

			if ec.Environment == args[0] {
				found = true
				env = ec
				break
			}
			envNames = append(envNames, ec.Environment)
		}

		if !found {
			fmt.Println(styles.ErrorStyle.Render("error:", "No such environment found. Valid oci environment that can be built are:", strings.Join(envNames, ", ")))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf oci <environment name>` to build the environment"))
			os.Exit(1)
		}

		// err = build.Build(env)
		// if err != nil {
		// 	fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		// 	os.Exit(1)
		// }

	},
}

func init() {
	OCICmd.Flags().StringVarP(&platform, "platform", "p", "", "The platform to build the image for")
}
