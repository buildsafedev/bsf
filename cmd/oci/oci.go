package oci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/build"
	binit "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/builddocker"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
	"github.com/buildsafedev/bsf/pkg/oci"
	"github.com/buildsafedev/bsf/pkg/platformutils"
)

var (
	platform, output string
	push, loadDocker bool
)
var (
	supportedPlatforms = []string{"linux/amd64", "linux/arm64"}
)

// OCICmd represents the export command
var OCICmd = &cobra.Command{
	Use:   "oci",
	Short: "Builds an OCI image",
	Long: `
	bsf oci <environment name> 
	bsf oci <environment name> --platform <platform>
	bsf oci <environment name> --platform <platform> --output <output directory>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// todo: we could provide a TUI list dropdown to select
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf oci <environment name>` to build an OCI image"))
			os.Exit(1)
		}

		env, p, err := ProcessPlatformAndConfig(platform, args[0])
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		platform = p

		sc, fh, err := binit.GetBSFInitializers()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = generate.Generate(fh, sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		if output == "" {
			output = "bsf-result"
		}

		err = bgit.Add("bsf/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = bgit.Ignore(output + "/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		symlink := "/result"

		err = nixcmd.Build(output+symlink, genOCIAttrName(env.Environment, platform))
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.HighlightStyle.Render("Generating artifacts..."))

		lockData, err := os.ReadFile("bsf.lock")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		lockFile := &hcl2nix.LockFile{}
		err = json.Unmarshal(lockData, lockFile)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		appDetails, graph, err := nixcmd.GetRuntimeClosureGraph(lockFile.App.Name, output, symlink)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		appDetails.Name = env.Name

		tos, tarch := platformutils.FindPlatform(platform)
		err = build.GenerateArtifcats(output, symlink, lockFile, appDetails, graph, tos, tarch)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Build completed successfully, please check the %s directory", output)))

		if loadDocker {
			fmt.Println(styles.HighlightStyle.Render("Loading image to docker daemon..."))

			expectedInstall := true
			currentContext, err := builddocker.GetCurrentContext()
			if err != nil {
				expectedInstall = false
			}
			contextEP, err := builddocker.ReadContextEndpoints()
			if err != nil {
				expectedInstall = false
			}
			if currentContext == "" {
				currentContext = "default"
			}
			if contextEP == nil {
				contextEP = make(map[string]string)
			}

			if _, ok := contextEP[currentContext]; !ok {
				contextEP[currentContext] = "unix:///var/run/docker.sock"
			}

			err = oci.LoadDocker(contextEP[currentContext], output+"/result", env.Name)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				if !expectedInstall {
					fmt.Println(styles.ErrorStyle.Render("error:", "Is Docker installed?"))
				}
				os.Exit(1)
			}

			fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s loaded to docker daemon", env.Name)))

		}

		if push {
			fmt.Println(styles.HighlightStyle.Render("Pushing image to registry..."))
			err = oci.Push(output+"/result", env.Name)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s pushed to registry", env.Name)))
		}

	},
}

// ProcessPlatformAndConfig processes the platform and config file
func ProcessPlatformAndConfig(plat string, envName string) (hcl2nix.OCIArtifact, string, error) {
	if plat == "" {
		tos, tarch := platformutils.FindPlatform(plat)
		plat = tos + "/" + tarch
	}

	pfound := false
	for _, sp := range supportedPlatforms {
		if strings.Contains(plat, sp) {
			pfound = true
			break
		}
	}
	if !pfound {
		return hcl2nix.OCIArtifact{}, "", fmt.Errorf("Platform %s is not supported. Supported platforms are %s", platform, strings.Join(supportedPlatforms, ", "))
	}
	data, err := os.ReadFile("bsf.hcl")
	if err != nil {
		return hcl2nix.OCIArtifact{}, "", fmt.Errorf("error: %s", err.Error())
	}

	var dstErr bytes.Buffer
	conf, err := hcl2nix.ReadConfig(data, &dstErr)
	if err != nil {
		return hcl2nix.OCIArtifact{}, "", fmt.Errorf(dstErr.String())
	}

	envNames := make([]string, 0, len(conf.OCIArtifact))
	var found bool
	env := hcl2nix.OCIArtifact{}
	for _, ec := range conf.OCIArtifact {
		errStr := ec.Validate(conf)
		if errStr != nil {
			return hcl2nix.OCIArtifact{}, "", fmt.Errorf("Config for export block %s is invalid\n Error: %s", ec.Name, *errStr)
		}

		if ec.Environment == envName {
			found = true
			env = ec
			break
		}
		envNames = append(envNames, ec.Environment)
	}

	if !found {
		return hcl2nix.OCIArtifact{}, "", fmt.Errorf("error: No such environment found. Valid oci environment that can be built are: %s", strings.Join(envNames, ", "))
	}

	return env, plat, nil
}

func genOCIAttrName(env, platform string) string {
	// .#ociImages.x86_64-linux.ociImage_caddy-as-dir
	tostarch := ""
	switch platform {
	case "linux/amd64":
		tostarch = "x86_64-linux"
	case "linux/arm64":
		tostarch = "aarch64-linux"
	}
	return fmt.Sprintf("bsf/.#ociImages.%s.ociImage_%s-as-dir", tostarch, env)
}

func init() {
	OCICmd.Flags().StringVarP(&platform, "platform", "p", "", "The platform to build the image for")
	OCICmd.Flags().StringVarP(&output, "output", "o", "", "location of the build artifacts generated")
	OCICmd.Flags().BoolVarP(&loadDocker, "load-docker", "", false, "Load the image into docker daemon")
	OCICmd.Flags().BoolVarP(&push, "push", "", false, "Push the image to the registry")

}
