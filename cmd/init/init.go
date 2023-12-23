package init

import (
	"context"
	"io/fs"
	"os"

	"github.com/buildsafedev/bsf/pkg/instrumentation/logging"
	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "init setups package management for the project",
	Long: `init setups package management for the project. It setups Nix files based on the language detected.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.SetupLogger()
		log.Info("initialising")
		ctx := logging.InjectLogger(context.Background(), log)
		files, err := createBsfDirectory(ctx)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		if len(files) != 0 {
			log.Info("Project has already been initialised")
			os.Exit(0)
		}

		// create the nix files
	},
}

func createBsfDirectory(ctx context.Context) ([]fs.DirEntry, error) {
	// check if the directory exists
	files, err := os.ReadDir("bsf")
	if err != nil {
		// check if the error is because the directory doesn't exist
		if os.IsNotExist(err) {
			err = os.Mkdir("bsf", 0755)
			if err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}

func createNixFiles(ctx context.Context) error {
	// create the nix files

	return nil
}
