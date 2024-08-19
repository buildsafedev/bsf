package oci

import (
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
	platform, output, tag, path, destcreds                string
	push, loadDocker, loadPodman, devDeps, dfSwap, digest bool
)
var (
	supportedPlatforms = []string{"linux/amd64", "linux/arm64"}
)

func init() {
	OCICmd.Flags().StringVarP(&platform, "platform", "p", "", "The platform to build the image for")
	OCICmd.Flags().StringVarP(&output, "output", "o", "", "location of the build artifacts generated")
	OCICmd.Flags().BoolVarP(&loadDocker, "load-docker", "", false, "Load the image into docker daemon")
	OCICmd.Flags().BoolVarP(&loadPodman, "load-podman", "", false, "Load the image into podman")
	OCICmd.Flags().BoolVarP(&push, "push", "", false, "Push the image to the registry")
	OCICmd.Flags().BoolVarP(&devDeps, "dev", "", false, "Build base image for Dev Dependencies")
	OCICmd.Flags().BoolVarP(&dfSwap, "df-swap", "", false, "Modify base images in Dockerfile")
	OCICmd.Flags().StringVarP(&tag, "tag", "t", "", "The tag that will be replaced with original tag in Dockerfile")
	OCICmd.Flags().StringVar(&path, "path", "", "The path to Dockerfile")
	OCICmd.Flags().BoolVar(&digest, "digest", false, "push image by digest")
	OCICmd.Flags().StringVar(&destcreds, "dest-creds", "", "Authenticate to the registry")
}

// OCICmd represents the export command
var OCICmd = &cobra.Command{
	Use:   "oci",
	Short: "Builds an OCI image",
	Long: `
	bsf oci <artifact> 
	bsf oci <artifact> --platform <platform>
	bsf oci <artifact> --push --digest
	bsf oci <artifact> --push --dest-creds user:password
	bsf oci <artifact> --platform <platform> --output <output directory>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// todo: we could provide a TUI list dropdown to select
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf oci <environment name>` to build an OCI image"))
			os.Exit(1)
		}

		conf, err := hcl2nix.ReadHclFile("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		artifact, p, err := ProcessPlatformAndConfig(conf, platform, args[0])
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		platform = p

		if tag != "" && !dfSwap {
			newName, err := getNewName(artifact, tag)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
				os.Exit(1)
			}
			artifact.Name = newName
		}

		if dfSwap {
			if tag != "" {
				if err = modifyDockerfileWithTag(path, tag); err != nil {
					fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
					os.Exit(1)
				}
				fmt.Println(styles.SucessStyle.Render("dockerfile succesfully updated with tag:", tag))
			} else {
				fmt.Println(styles.HintStyle.Render("hint:", "use --tag flag to define a tag"))
			}
			os.Exit(1)
		}

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

		err = nixcmd.Build(output+symlink, genOCIAttrName(artifact.Artifact, platform, artifact.IsBase))
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
		appDetails.Name = artifact.Name

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

			err = oci.LoadDocker(contextEP[currentContext], output+"/result", artifact.Name)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				if !expectedInstall {
					fmt.Println(styles.ErrorStyle.Render("error:", "Is Docker installed?"))
				}
				os.Exit(1)
			}

			fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s loaded to docker daemon", artifact.Name)))

		}

		if loadPodman {
			fmt.Println(styles.HighlightStyle.Render("Loading image to podman..."))
			err = oci.LoadPodman(output+"/result", artifact.Name)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s loaded to podman", artifact.Name)))
		}

		if push {
			fmt.Println(styles.HighlightStyle.Render("Pushing image to registry..."))
			if digest {
				err = oci.Push(output+"/result", artifact.Name, destcreds, output+"/digest")
				if err != nil {
					fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
					os.Exit(1)
				}

				fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s pushed to registry", artifact.Name)))
				digest, err := os.ReadFile(output + "/digest")
				if err != nil {
					fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
					os.Exit(1)
				}
				fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Digest: %s", string(digest))))
			} else {
				err = oci.Push(output+"/result", artifact.Name, destcreds, "")
				if err != nil {
					fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
					os.Exit(1)
				}
				fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Image %s pushed to registry", artifact.Name)))
			}
		}
	},
}

// ProcessPlatformAndConfig processes the platform and config file
func ProcessPlatformAndConfig(conf *hcl2nix.Config, plat string, envName string) (hcl2nix.OCIArtifact, string, error) {
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

	envNames := make([]string, 0, len(conf.OCIArtifact))
	var found bool
	artifact := hcl2nix.OCIArtifact{}
	for _, ec := range conf.OCIArtifact {
		errStr := ec.Validate(conf)
		if errStr != nil {
			return hcl2nix.OCIArtifact{}, "", fmt.Errorf("Config for oci block %s is invalid\n Error: %s", ec.Name, *errStr)
		}

		if ec.Artifact == envName {
			found = true
			artifact = ec
			break
		}
		envNames = append(envNames, ec.Artifact)
	}

	if !found {
		return hcl2nix.OCIArtifact{}, "", fmt.Errorf("error: No such environment found. Valid oci environment that can be built are: %s", strings.Join(envNames, ", "))
	}

	return artifact, plat, nil
}

func modifyDockerfileWithTag(path, tag string) error {
	var dockerfilePath string
	if path != "" {
		dockerfilePath = path + "/Dockerfile"
	} else {
		dockerfilePath = "./Dockerfile"
	}

	file, err := os.Open(dockerfilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	resLines, err := builddocker.ModifyDockerfile(file, devDeps, tag)
	if err != nil {
		return err
	}

	err = os.WriteFile(dockerfilePath, []byte(strings.Join(resLines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getNewName(artifact hcl2nix.OCIArtifact, tag string) (string, error) {
	var newName string
	if strings.Contains(artifact.Name, ":") {
		parts := strings.Split(artifact.Name, ":")
		if len(parts) > 0 {
			newName = fmt.Sprintf("%s:%s", parts[0], tag)
		} else {
			newName = fmt.Sprintf("%s:%s", artifact.Name, tag)
		}
	} else {
		newName = fmt.Sprintf("%s:%s", artifact.Name, tag)
	}
	return newName, nil
}

func genOCIAttrName(env, platform string, isBase bool) string {
	var arch string

	switch platform {
	case "linux/amd64":
		arch = "x86_64-linux"
	case "linux/arm64":
		arch = "aarch64-linux"
	default:
		arch = "unknown"
	}

	base := fmt.Sprintf("bsf/.#ociImage_%s.%s.ociImage_%s_", env, arch, env)

	if isBase {
		return base + "base-as-dir"
	}

	return base + "app-as-dir"
}
