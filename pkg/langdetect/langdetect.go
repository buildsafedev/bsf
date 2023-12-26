package langdetect

import (
	"os"
	"path/filepath"
)

type (
	// ProjectType is the type of project
	ProjectType string
)

const (
	// GoModule is the project type for Go modules
	GoModule ProjectType = "GoModule"
	// Unknown is the project type for unknown project types
	Unknown ProjectType = "Unknown"
)

// FindProjectType detects the programming language/package manager of the current project.
func FindProjectType() (ProjectType, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// List files in the current directory
	files, err := filepath.Glob(filepath.Join(currentDir, "*"))
	if err != nil {
		return "", err
	}

	// Check for specific project files
	for _, file := range files {
		switch filepath.Base(file) {
		case "go.mod":
			return GoModule, nil
		}
	}

	return Unknown, nil
}
