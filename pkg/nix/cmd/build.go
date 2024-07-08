package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// Build invokes nix build to build the project
func Build(dir string, attribute string) error {
	if attribute == "" {
		attribute = "bsf/."
	}
	cmd := exec.Command("nix", "build", attribute, "-o", dir)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("could not get stderr pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("could not get stdout pipe: %v", err)
	}
	go func() {
		merged := io.MultiReader(stderr, stdout)
		scanner := bufio.NewScanner(merged)
		for scanner.Scan() {
			msg := scanner.Text()
			dir, err:= os.Getwd()
			if err!=nil{
				panic(err)
			}
			warning:= fmt.Sprintf("warning: Git tree '%s' is dirty", dir)
			if msg==warning{
				msg = fmt.Sprintf("warning: Git tree '%s' is dirty.\nThis implies you have not checked-in files in the git work tree (hint: git add)", dir)
			}
			fmt.Printf("%s\n", msg)
		}
	}()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}
	return nil
}