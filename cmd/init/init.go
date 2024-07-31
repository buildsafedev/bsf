package init

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/precheck"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "init setups package management for the project",
	Long: `init setups package management for the project. It setups Nix files based on the language detected.
	`,

	PreRun: func(cmd *cobra.Command, args []string) {
		precheck.AllPrechecks()
	},

	Run: func(cmd *cobra.Command, args []string) {
		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		isBaseImage, err := YesNoPrompt("Do you want to build a base image?")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		var imageName string
		var pt langdetect.ProjectType
		var pd *langdetect.ProjectDetails

		if isBaseImage {
			imageName, err = IoPrompt("What should the image name be?")
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			pt = langdetect.BaseImage

			isAddDeps, err := YesNoPrompt("Would you like to add common Go dependencies (compiler, linter, debugger, etc)?")
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}

			if isAddDeps {
				// Create configuration with Go dependencies
				config := genGoModuleConf(&langdetect.ProjectDetails{Name: imageName})

				// Write the configuration to the file
				err = hcl2nix.WriteConfig(config, os.Stdout)
				if err != nil {
					fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
					os.Exit(1)
				}

				fmt.Println(styles.HighlightStyle.Render("Base image configuration with common Go dependencies has been generated."))
			}
		}

		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		config, err := generatehcl2NixConf(pt, pd, imageName)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		err = hcl2nix.WriteConfig(config, os.Stdout)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		m := model{sc: sc, pt: pt, baseImgName: imageName}
		m.resetSpinner()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			os.Exit(1)
		}
	},
}

// YesNoPrompt is used to promt a user for a bool based answer
func YesNoPrompt(label string) (bool, error) {
	choices := "Y/n"

	r := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fprintf(os.Stderr, styles.HighlightStyle.Render("%s (%s) "), strings.TrimSpace(label), choices)
	s, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "y" || s == "yes" {
		return true, nil
	}
	if s == "n" || s == "no" {
		return false, nil
	}
	return false, nil
}

// IoPrompt is used to promt a user for a string based answer
func IoPrompt(label string) (string, error) {
	r := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fprintf(os.Stderr, styles.HighlightStyle.Render("%s : "), strings.TrimSpace(label))
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("please define a proper name for the image")
	}
	return s, nil
}

// GetBSFInitializers generates the nix files
func GetBSFInitializers() (bsfv1.SearchServiceClient, *hcl2nix.FileHandlers, error) {
	if _, err := os.Stat("bsf.hcl"); err != nil {
		fmt.Println(styles.HintStyle.Render("hint: ", "run `bsf init` to initialize the project"))
		return nil, nil, fmt.Errorf("error: %s\nHas the project been initialized?", err.Error())
	}

	conf, err := configure.PreCheckConf()
	if err != nil {
		return nil, nil, fmt.Errorf("error: %s", err.Error())
	}

	fh, err := hcl2nix.NewFileHandlers(true)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %s", err.Error())
	}

	sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %s", err.Error())
	}

	return sc, fh, nil
}

// CleanUp removes the bsf config if any error occurs in init process (ctrl+c or any init process stage)
func cleanUp() {
	configs := []string{"bsf", "bsf.hcl", "bsf.lock"}

	for _, f := range configs {
		os.RemoveAll(f)
	}
}
