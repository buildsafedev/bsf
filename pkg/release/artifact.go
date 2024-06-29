package release

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// WalkDir recursively walks through the directory and notes down all the file names.
func WalkDir(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			subDirFiles, err := WalkDir(fullPath)
			if err != nil {
				return nil, err
			}
			files = append(files, subDirFiles...)
		} else {
			files = append(files, fullPath)
		}
	}

	return files, nil
}

// CreateTarGzFromSymlink creates a tar.gz file from the specified directory.
func CreateTarGzFromSymlink(sourceDir, targetFile string) error {
	tarFile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzWriter := gzip.NewWriter(tarFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return addSymlinkToTar(sourceDir, tarWriter, path)
		}

		return addFileToTar(sourceDir, tarWriter, path, info)
	})
}

func addSymlinkToTar(sourceDir string, tarWriter *tar.Writer, path string) error {
	linkTarget, err := os.Readlink(path)
	if err != nil {
		return err
	}

	// Get FileInfo of symlink target
	targetInfo, err := os.Lstat(linkTarget)
	if err != nil {
		return err
	}

	if targetInfo.IsDir() {
		return filepath.Walk(linkTarget, func(subpath string, subinfo os.FileInfo, suberr error) error {
			if suberr != nil {
				return suberr
			}

			return addFileToTar(sourceDir, tarWriter, subpath, subinfo)
		})
	}

	// If the symlink points to a file, treat it as a regular file
	return addFileToTar(sourceDir, tarWriter, linkTarget, targetInfo)
}

func addFileToTar(sourceDir string, tarWriter *tar.Writer, path string, info os.FileInfo) error {
	// Resolve sourceDir if it's a symlink
	resolvedSourceDir, err := resolveSymlink(sourceDir)
	if err != nil {
		return err
	}

	// Calculate the relative path based on the resolved sourceDir
	relPath, err := calculateRelativePath(resolvedSourceDir, path)
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filepath.ToSlash(relPath)

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tarWriter, file)
	return err
}

func resolveSymlink(path string) (string, error) {
	// Resolve symlink target
	target, err := os.Readlink(path)
	if err != nil {
		return "", err
	}
	return target, nil
}

func calculateRelativePath(sourceDir, path string) (string, error) {
	// If sourceDir is not a prefix of path, return an error
	if !strings.HasPrefix(path, sourceDir) {
		return "", fmt.Errorf("path %s is not within source directory %s", path, sourceDir)
	}

	// Calculate relative path based on sourceDir
	relPath, err := filepath.Rel(sourceDir, path)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(relPath), nil
}

// CreateTarGzFromDir creates a .tar.gz archive from a directory.
func CreateTarGzFromDir(sourceDir, targzPath string) error {
	// Step 1: Open the target .tar.gz file for writing
	targzFile, err := os.Create(targzPath)
	if err != nil {
		return err
	}
	defer targzFile.Close()

	gzWriter := gzip.NewWriter(targzFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	err = filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, file)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
