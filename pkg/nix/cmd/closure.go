package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/awalterschulze/gographviz"
)

// App represents the application
type App struct {
	Name    string
	Version string
	Hash    string
}

// GetRuntimeClosureGraph returns the runtime closure graph for the project
func GetRuntimeClosureGraph() (*App, *gographviz.Graph, error) {
	app, err := GetAppDetails()
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command("nix-store", "-q", "--graph", "result")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("failed with %s", cmd.Stderr)
	}

	graphAst, err := gographviz.ParseString(stdout.String())
	if err != nil {

		return nil, nil, fmt.Errorf("failed to parse graph: %s", err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, nil, fmt.Errorf("failed to analyse graph: %s", err)
	}

	return app, graph, nil
}

// GetAppDetails checks if the result symlink exists
func GetAppDetails() (*App, error) {
	target, err := os.Readlink("result")
	if err != nil {
		return nil, fmt.Errorf("failed to read symlink: %v", err)
	}

	app, err := parseAppDetails(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse app details: %v", err)
	}

	return app, nil
}

func parseAppDetails(path string) (*App, error) {
	path = strings.TrimPrefix(path, "/nix/store/")
	parts := strings.Split(path, "-")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	return &App{
		Hash:    parts[0],
		Name:    parts[len(parts)-2],
		Version: parts[len(parts)-1],
	}, nil
}
