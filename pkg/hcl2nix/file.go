package hcl2nix

import (
	"fmt"
	"io/fs"
	"os"
)

// FileHandlers holds file handlers
type FileHandlers struct {
	ModFile      *os.File
	LockFile     *os.File
	FlakeFile    *os.File
	DefFlakeFile *os.File
}

// NewFileHandlers creates new file handlers
func NewFileHandlers(expectInit bool) (*FileHandlers, error) {
	var err error
	_, err = createBsfDirectory()
	if err != nil {
		return nil, err
	}

	modFile, err := CreateModFile(expectInit)
	if err != nil {
		return nil, err
	}

	lockFile, err := os.Create("bsf.lock")
	if err != nil {
		return nil, err
	}

	flakeFile, err := os.Create("bsf/flake.nix")
	if err != nil {
		return nil, err
	}

	defFlakeFile, err := os.Create("bsf/default.nix")
	if err != nil {
		return nil, err
	}

	return &FileHandlers{
		ModFile:      modFile,
		LockFile:     lockFile,
		FlakeFile:    flakeFile,
		DefFlakeFile: defFlakeFile,
	}, nil
}

// CreateModFile creates bsf.hcl file
func CreateModFile(expectInit bool) (*os.File, error) {
	bsfHcl := "bsf.hcl"
	var modFile *os.File
	var exists bool
	if _, err := os.Stat(bsfHcl); os.IsNotExist(err) {
		if expectInit {
			return nil, fmt.Errorf("Project not initialised. bsf.hcl not found")
		}
		modFile, err = os.Create(bsfHcl)
		if err != nil {
			return nil, err
		}
	} else {
		exists = true
		modFile, err = os.Open(bsfHcl)
		if err != nil {
			return nil, err
		}
	}

	if exists != expectInit {
		return nil, fmt.Errorf("Project already initialised. bsf.hcl found")
	}

	return modFile, nil
}

// GetOrCreateFile gets or creates a file if it doesn't exist
func GetOrCreateFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.Open(path)
}

func createBsfDirectory() ([]fs.DirEntry, error) {
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
