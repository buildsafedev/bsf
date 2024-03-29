package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Add adds the path to the git work tree
func Add(path string) error {
	var r *git.Repository
	var err error
	r, err = git.PlainOpen(".")
	if err == git.ErrRepositoryNotExists {
		// If it's not a Git repository, initialize it
		r, err = git.PlainInit(".", false)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Add all changes to the working directory
	_, err = w.Add(path)
	if err != nil {
		return err
	}
	return nil
}

// Ignore adds the path to the .gitignore file
func Ignore(path string) error {
	gitignorePath := filepath.Join(".", ".gitignore")

	// Check if .gitignore exists
	_, err := os.Stat(gitignorePath)
	if os.IsNotExist(err) {
		// Create .gitignore if it doesn't exist
		_, err := os.Create(gitignorePath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Read .gitignore
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return err
	}

	// Check if path is already in .gitignore
	if strings.Contains(string(content), path) {
		return nil
	}

	// Append path to .gitignore content
	content = append(content, []byte("\n"+path)...)

	// Write new content to .gitignore
	err = os.WriteFile(gitignorePath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
