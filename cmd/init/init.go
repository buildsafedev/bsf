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

		res := YesNoPrompt("Do you want docker file to be initialized", false)

		var pd langdetect.ProjectDetails
		var pt langdetect.ProjectType
		if res {
			pd.Name = *IOprompt("What should the name be?", "").(*string)
		} else {
			var projName = IOprompt("What type of projects you are setting up ? (Go/Rust/JavaScript/Poetry)", nil)
			fmt.Println("TYPE : ", projName)
			var lang langdetect.ProjectType
			switch projName.(string) {
			case "Go":
				lang = langdetect.GoModule
			case "Rust":
				lang = langdetect.RustCargo
			case "JS":
				lang = langdetect.JsNpm
			case "Poetry":
				lang = langdetect.PythonPoetry
			}
			pt = lang
		}

		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			os.Exit(1)
		}

		m := model{sc: sc, pd: &pd, pt: pt}
		m.resetSpinner()
		if _, err := tea.NewProgram(m).Run(); err != nil {
			os.Exit(1)
		}
	},
}

func YesNoPrompt(label string, defaultChoice bool) bool {
	choices := "Y/n"
	if !defaultChoice {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", strings.TrimSpace(label), choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return defaultChoice
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func IOprompt(label string, defaultvalue any) any {
	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Fprintf(os.Stderr, "%s : ", strings.TrimSpace(label))
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}

		return &s
	}
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
