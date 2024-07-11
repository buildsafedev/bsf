package cmd

import (
	"bufio"
	"errors"
	"fmt"
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
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdErrChan:= make(chan string)

	go func(){
		scanner:= bufio.NewScanner(stderr)
		for scanner.Scan() {
			stdErr:= scanner.Text()
			dir, err:= os.Getwd()
			if err!=nil{
				stdErrChan <- err.Error()
			}
			warning:= fmt.Sprintf("warning: Git tree '%s' is dirty", dir)
			if stdErr==warning{
				stdErr = fmt.Sprintf("warning: Git tree '%s' is dirty.\nThis implies you have not checked-in files in the git work tree (hint: git add)", dir)
			}
			stdErr = fmt.Sprint(stdErr, "\n")
			os.Stderr.Write([]byte(stdErr))
		}
	}()
	
	go func(){
		scanner:= bufio.NewScanner(stdout)
		for scanner.Scan() {
			stdOut:= scanner.Text()
			stdOut = fmt.Sprint(stdOut, "\n")
			os.Stdout.Write([]byte(stdOut))
		}
	}()
	
	select{
	case stdErr:=<-stdErrChan:
		return errors.New(stdErr)
	default:
		break
	}
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}
	return nil
}