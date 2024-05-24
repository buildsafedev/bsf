package init

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/precheck"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
	"github.com/buildsafedev/bsf/pkg/nix/cmd"
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

		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		err = initializeProject(sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			cleanUp()
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render("Initialized successfully!"))
	},
}

func initializeProject(sc bsfv1.SearchServiceClient) error {
	fmt.Println(styles.TextStyle.Render("Initializing project, detecting project language..."))

	pt, pd, err := langdetect.FindProjectType()
	if err != nil {
		return err
	}

	fmt.Println(styles.TextStyle.Render("Detected language as " + string(pt)))
	if pt == langdetect.Unknown {
		return fmt.Errorf("project language isn't currently supported, some features might not work")
	}

	fmt.Println(styles.TextStyle.Render("Resolving dependencies..."))

	fh, err := hcl2nix.NewFileHandlers(false)
	if err != nil {
		return err
	}
	defer fh.ModFile.Close()
	defer fh.LockFile.Close()
	defer fh.FlakeFile.Close()
	defer fh.DefFlakeFile.Close()

	conf, err := generatehcl2NixConf(pt, pd)
	if err != nil {
		return err
	}

	err = hcl2nix.WriteConfig(conf, fh.ModFile)
	if err != nil {
		return err
	}

	err = generate.Generate(fh, sc)
	if err != nil {
		return err
	}

	err = cmd.Lock()
	if err != nil {
		return err
	}

	err = bgit.Add("bsf/")
	if err != nil {
		return err
	}

	err = bgit.Ignore("bsf-result/")
	if err != nil {
		return err
	}

	return nil
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
