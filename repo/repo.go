package repo

import (
	"os"
	"patchy/util"
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
		hasRepoDir, err := util.DoesFileExist(filepath.Join(dir, ".patchy"))
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
	return "", &ErrNotInRepo{}
}

func FindRepoRoot() (string, error) {
	repoDir, err := FindRepoDir()
	if err != nil {
		return "", err
	}
	return filepath.Dir(repoDir), nil
}

func IsFileInRepo(path string) (bool, error) {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		return false, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	if _, err := filepath.Rel(repoRoot, absPath); err != nil {
		return false, nil
	}
	return true, nil
}
