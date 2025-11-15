package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitRepo(path string) (string, error) {
	if _, e := FindRepoDir(); e == nil {
		return "", fmt.Errorf("InitRepo: %w", ErrAlreadyInRepo)
	}

	repoPath, err := filepath.Abs(filepath.Join(path, ".patchy"))
	if err != nil {
		return "", fmt.Errorf("InitRepo: %w", err)
	}

	if err = os.MkdirAll(filepath.Join(repoPath, "objects"), os.ModePerm); err != nil {
		_ = os.RemoveAll(repoPath)
		return "", fmt.Errorf("InitRepo: %w", err)
	}
	if err = os.MkdirAll(filepath.Join(repoPath, "refs", "heads"), os.ModePerm); err != nil {
		_ = os.RemoveAll(repoPath)
		return "", fmt.Errorf("InitRepo: %w", err)
	}
	if err = os.MkdirAll(filepath.Join(repoPath, "refs", "tags"), os.ModePerm); err != nil {
		_ = os.RemoveAll(repoPath)
		return "", fmt.Errorf("InitRepo: %w", err)
	}

	if err = os.WriteFile(filepath.Join(repoPath, "HEAD"), []byte("ref: refs/heads/main"), 0666); err != nil {
		_ = os.RemoveAll(repoPath)
		return "", fmt.Errorf("InitRepo: %w", err)
	}

	return repoPath, nil
}
