package dockerfile

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"
	"github.com/stacklok/frizbee/pkg/replacer"
	"github.com/stacklok/frizbee/pkg/utils/config"
)

var DGCmd = &cobra.Command{
	Use:     "digests",
	Short:   "Replace Dockerfile image tags with immutable digests",
	Aliases: []string{"dg"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf dockerfile digests <Dockerfile>` to replace image tags with digests"))
			os.Exit(1)
		}

		dockerfile := args[0]

		r := replacer.NewContainerImagesReplacer(config.DefaultConfig())

		str, err := r.ParsePath(context.TODO(), dockerfile)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error in parsing Dockerfile contents", err.Error()))
			os.Exit(1)
		}

		if err = processOutput(dockerfile, str.Processed, str.Modified); err != nil {
			fmt.Println(styles.ErrorStyle.Render("error in writing Dockerfile contents", err.Error()))
			os.Exit(1)
		}
	},
}

func processOutput(path string, processed []string, modified map[string]string) error {
	basedir := filepath.Dir(path)
	bfs := osfs.New(basedir, osfs.WithBoundOS())
	var out io.Writer

	for path, content := range modified {
		f, err := bfs.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}

		defer func(f billy.File) {
			if err := f.Close(); err != nil {
				fmt.Println(styles.ErrorStyle.Render("failed to close file %s: %v", path, err.Error()))
			}
		}(f)

		out = f

		_, err = fmt.Fprintf(out, "%s", content)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", path, err)
		}
	}

	return nil
}
