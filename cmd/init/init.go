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

		isBaseImage := yesNoPrompt("Do you want to build a base image?")
		var imageName string
		var pt langdetect.ProjectType
		if isBaseImage {
			imageName = *ioPrompt("What should the image name be?").(*string)
			pt = langdetect.BaseImage
		}

		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			os.Exit(1)
		}

		m := model{sc: sc, pt: pt, baseImgName: imageName}
		m.resetSpinner()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			os.Exit(1)
		}
	},
}

func yesNoPrompt(label string) bool {
	choices := "Y/n"

	r := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fprintf(os.Stderr, "%s (%s) ", strings.TrimSpace(label), choices)
	s, _ = r.ReadString('\n')
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s == "y" || s == "yes" {
		return true
	}
	if s == "n" || s == "no" {
		return false
	}
	return false
}

func ioPrompt(label string) any {
	r := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fprintf(os.Stderr, "%s : ", strings.TrimSpace(label))
	s, _ = r.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
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
