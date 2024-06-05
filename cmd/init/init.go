package init

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
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
	// Set up the progress writer
	pw, steps := setupProgressTracker()

	trackers := make([]*progress.Tracker, len(steps))
	for i, step := range steps {
		trackers[i] = &progress.Tracker{Message: step, Total: 100}
		pw.AppendTracker(trackers[i])
	}

	go pw.Render()

	// Initialize progress tracking
	updateProgress := func(tracker *progress.Tracker, progress int64) {
		tracker.SetValue(progress)
		if progress >= 100 {
			tracker.MarkAsDone()
		}
	}

	// Detect project language
	pt, pd, err := langdetect.FindProjectType()
	if err != nil {
		return err
	}
	updateProgress(trackers[0], 100)

	// Resolve dependencies
	fh, err := hcl2nix.NewFileHandlers(false)
	if err != nil {
		return err
	}
	defer fh.ModFile.Close()
	defer fh.LockFile.Close()
	defer fh.FlakeFile.Close()
	defer fh.DefFlakeFile.Close()
	updateProgress(trackers[1], 100)

	// Write configuration
	conf, err := generatehcl2NixConf(pt, pd)
	if err != nil {
		return err
	}
	err = hcl2nix.WriteConfig(conf, fh.ModFile)
	if err != nil {
		return err
	}
	updateProgress(trackers[2], 100)

	// Generate files
	err = generate.Generate(fh, sc)
	if err != nil {
		return err
	}
	updateProgress(trackers[3], 100)

	// Lock dependencies
	err = cmd.Lock(func(progress int) {
		updateProgress(trackers[4], int64(progress))
	})
	if err != nil {
		return err
	}

	// Add to git
	err = bgit.Add("bsf/")
	if err != nil {
		return err
	}

	// Set up git ignore
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

func setupProgressTracker() (progress.Writer, []string) {
	pw := progress.NewWriter()
	pw.SetAutoStop(true)
	pw.SetTrackerLength(25)
	pw.SetMessageLength(50)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	pw.Style().Visibility.ETA = true
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Time = true

	// Define the steps and create trackers for each
	steps := []string{
		"Detecting project language...",
		"Resolving dependencies...",
		"Writing configuration...",
		"Generating files...",
		"Locking dependencies...",
	}

	return pw, steps
}
