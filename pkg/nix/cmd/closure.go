package cmd

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/awalterschulze/gographviz"
)

// GetRuntimeClosureGraph returns the runtime closure graph for the project
func GetRuntimeClosureGraph() (*gographviz.Graph, error) {
	cmd := exec.Command("nix-store", "-q", "--graph", "result")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed with %s", cmd.Stderr)
	}

	graphAst, err := gographviz.ParseString(stdout.String())
	if err != nil {

		return nil, fmt.Errorf("failed to parse graph: %s", err)
	}

	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(graphAst, graph); err != nil {
		return nil, fmt.Errorf("failed to analyse graph: %s", err)
	}

	return graph, nil
}
