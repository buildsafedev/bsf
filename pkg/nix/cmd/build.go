package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

// Build invokes nix build to build the project
func Build(dir string, attribute string) error {
	if attribute == "" {
		attribute = "bsf/."
	}
	cmd := exec.Command("nix", "build", attribute, "-o", dir)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	g:= new(errgroup.Group)

	g.Go(func()error{
		return ManageStdErr(stderr)
	})

	g.Go(func()error{
		return ManageStdOutput(stdout)
	})

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}

	return g.Wait()
}

func ManageStdErr(stderr io.ReadCloser) error {
	scanner:= bufio.NewScanner(stderr)
		for scanner.Scan() {
			stdErr:= scanner.Text()
			dir, err:= os.Getwd()
			if err!=nil{
				return err
			}
			warning:= fmt.Sprintf("warning: Git tree '%s' is dirty", dir)
			if stdErr==warning{
				stdErr = fmt.Sprintf("warning: Git tree '%s' is dirty.\nThis implies you have not checked-in files in the git work tree (hint: git add)", dir)
			}
			stdErr = fmt.Sprint(stdErr, "\n")
			os.Stderr.Write([]byte(stdErr))
		}
	return nil
}

func ManageStdOutput(stdout io.ReadCloser) error {
	scanner:= bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdOut:= scanner.Text()
			stdOut = fmt.Sprint(stdOut, "\n")
			os.Stdout.Write([]byte(stdOut))
		}
	return nil
}