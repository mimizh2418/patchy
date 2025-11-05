package util

import (
	"errors"
	"os"
	"path/filepath"
)

var foundRepoDir = false
var repoDir string

func FindRepoDir() (string, error) {
	if foundRepoDir {
		return repoDir, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir != filepath.Dir(dir) {
		hasRepoDir, err := DoesFileExist(filepath.Join(dir, ".patchy"))
		if err != nil {
			return "", err
		}

		if hasRepoDir {
			foundRepoDir = true
			repoDir = filepath.Join(dir, ".patchy")
			return repoDir, nil
		}
		dir = filepath.Dir(dir)
	}
	return "", errors.New("current directory is not part of a repository")
}

func FindRepoRoot() (string, error) {
	repoDir, err := FindRepoDir()
	if err != nil {
		return "", err
	}
	return filepath.Dir(repoDir), nil
}
