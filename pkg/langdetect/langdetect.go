package langdetect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

type (
	// ProjectType is the type of project
	ProjectType string
)

const (
	// GoModule is the project type for Go modules
	GoModule ProjectType = "GoModule"

	// PythonPoetry is the project type for Python Poetry projects
	PythonPoetry ProjectType = "PythonPoetry"
	// Unknown is the project type for unknown project types
	Unknown ProjectType = "Unknown"
)

// ProjectDetails holds the details of the current project
type ProjectDetails struct {
	Entrypoint string
	Name       string
}

var supportedLanguages = []string{string(GoModule), string(PythonPoetry)}

// FindProjectType detects the programming language/package manager of the current project.
func FindProjectType() (ProjectType, *ProjectDetails, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", nil, err
	}

	// List files in the current directory
	files, err := filepath.Glob(filepath.Join(currentDir, "*"))
	if err != nil {
		return "", nil, err
	}

	// Check for specific project files
	for _, file := range files {
		switch filepath.Base(file) {
		case "go.mod":
			pd, err := pdFromGoMod(file)
			return GoModule, pd, err

		case "poetry.lock":
			return PythonPoetry, &ProjectDetails{}, nil

		default:
			err = fmt.Errorf("unable to detect the language ,supported languages: " + (strings.Join(supportedLanguages, ",") + "."))
		}
	}

	return Unknown, nil, err
}

func pdFromGoMod(goModPath string) (*ProjectDetails, error) {
	// read the file and check if it has a module name
	f, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.Parse(goModPath, f, nil)
	if err != nil {
		return nil, err
	}

	binaryName := binaryFromModule(mod)

	return &ProjectDetails{
		Name: binaryName,
	}, nil
}

func binaryFromModule(mod *modfile.File) string {
	pathParts := strings.Split(mod.Module.Mod.Path, "/")
	lastPart := pathParts[len(pathParts)-1]

	// If the last part is a version string, return the second last part
	if strings.HasPrefix(lastPart, "v") && len(lastPart) > 1 && strings.Trim(lastPart[1:], "0123456789") == "" {
		return pathParts[len(pathParts)-2]
	}

	// Otherwise, return the last part
	return lastPart
}
