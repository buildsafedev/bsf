package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildsafedev/bsf/pkg/langdetect"
	"github.com/go-git/go-git/v5"
)
type ErrFileNotAddedToVersionControl struct{
	fileName string
}
func (e *ErrFileNotAddedToVersionControl) Error() string {
	return fmt.Sprint(e.fileName, " is not added to version control")
}

// Add adds the path to the git work tree
func Add(path string) error {
	var r *git.Repository
	var err error

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	r, err = git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
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

	root := w.Filesystem.Root()
	leafDir := getLeafDir(root+"/", cwd)
	if leafDir != "" {
		path = leafDir + "/" + path
	}
	// Add all changes to the working directory
	err = w.AddWithOptions(&git.AddOptions{
		Path: path,
	})
	if err != nil {
		return err
	}
	status, err :=w.Status()
	if err !=nil{
		return err
	}
	pt, _, err:=langdetect.FindProjectType()
	if err!=nil{
		return err
	}
	entryFile:= langdetect.GetEntryFileOfProject(pt)
	_, exis:= status[entryFile]
	// We need to assume that if file is not in status map, then it is comitted
	if exis{
		fl:= status.File(entryFile)
		if fl.Staging==git.Untracked{
			return &ErrFileNotAddedToVersionControl{
				fileName: entryFile,
			}
		}
	}
	return nil
}

func getLeafDir(root string, path string) string {
	if strings.Compare(root, path+"/") == 0 {
		return ""
	}
	return strings.TrimPrefix(path, root)
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
