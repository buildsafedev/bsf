package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/spf13/cobra"
)

// Helper is a common struct for implementing a CLI command that replaces
// files.
type Helper struct {
	DryRun        bool
	Quiet         bool
	ErrOnModified bool
	Regex         string
	Cmd           *cobra.Command
}

// NewHelper creates a new CLI Helper struct.
func NewHelper(cmd *cobra.Command) (*Helper, error) {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return nil, fmt.Errorf("failed to get dry-run flag: %w", err)
	}
	errOnModified, err := cmd.Flags().GetBool("error")
	if err != nil {
		return nil, fmt.Errorf("failed to get error flag: %w", err)
	}
	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		return nil, fmt.Errorf("failed to get quiet flag: %w", err)
	}
	regex, err := cmd.Flags().GetString("regex")
	if err != nil {
		return nil, fmt.Errorf("failed to get regex flag: %w", err)
	}

	return &Helper{
		Cmd:           cmd,
		DryRun:        dryRun,
		ErrOnModified: errOnModified,
		Quiet:         quiet,
		Regex:         regex,
	}, nil
}

// DeclareFrizbeeFlags declares the flags common to all replacer commands.
func DeclareFrizbeeFlags(cmd *cobra.Command, enableOutput bool) {
	cmd.Flags().BoolP("dry-run", "n", false, "don't modify files")
	cmd.Flags().BoolP("quiet", "q", false, "don't print anything")
	cmd.Flags().BoolP("error", "e", false, "exit with error code if any file is modified")
	cmd.Flags().StringP("regex", "r", "", "regex to match artifact references")
	cmd.Flags().StringP("platform", "p", "", "platform to match artifact references, e.g. linux/amd64")
	if enableOutput {
		cmd.Flags().StringP("output", "o", "table", "output format. Can be 'json' or 'table'")
	}
}

// Logf logs the given message to the given command's stderr if the command is
// not quiet.
func (r *Helper) Logf(format string, args ...interface{}) {
	if !r.Quiet {
		fmt.Fprintf(r.Cmd.ErrOrStderr(), format, args...) // nolint:errcheck
	}
}

// ProcessOutput processes the given output files.
// If the command is quiet, the output is discarded.
// If the command is a dry run, the output is written to the command's stdout.
// Otherwise, the output is written to the given filesystem.
func (r *Helper) ProcessOutput(path string, processed []string, modified map[string]string) error {
	basedir := filepath.Dir(path)
	bfs := osfs.New(basedir, osfs.WithBoundOS())
	var out io.Writer
	for _, path := range processed {
		if !r.Quiet {
			r.Logf("Processed: %s\n", path)
		}
	}
	for path, content := range modified {
		if r.Quiet {
			out = io.Discard
		} else if r.DryRun {
			out = r.Cmd.OutOrStdout()
		} else {
			f, err := bfs.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}

			defer func() {
				if err := f.Close(); err != nil {
					fmt.Fprintf(r.Cmd.ErrOrStderr(), "failed to close file %s: %v", path, err) // nolint:errcheck
				}
			}()

			out = f
		}
		if !r.Quiet {
			r.Logf("Modified: %s\n", path)
		}
		_, err := fmt.Fprintf(out, "%s", content)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", path, err)
		}
	}

	return nil
}

// IsPath returns true if the given path is a file or directory.
func IsPath(pathOrRef string) bool {
	_, err := os.Stat(pathOrRef)
	return err == nil
}
