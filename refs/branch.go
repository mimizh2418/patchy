package refs

import (
	"fmt"
	"io/fs"
	"os"
	"patchy/repo"
	"path/filepath"
)

type Branch struct {
	Name       string
	CommitHash string
}

func NewBranch(name string, revSpec string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return fmt.Errorf("NewBranch: %w", err)
	}
	branches, err := ListBranches()
	if err != nil {
		return fmt.Errorf("NewBranch: %w", err)
	}
	for _, branch := range branches {
		if branch.Name == name {
			return fmt.Errorf("NewBranch: branch %s already exists", name)
		}
	}
	commitHash, err := ParseRev(revSpec)
	if err != nil {
		return fmt.Errorf("NewBranch: %w", err)
	}
	branchPath := filepath.Join(repoDir, "refs", "heads", name)
	branchDir := filepath.Dir(branchPath)
	if err := os.MkdirAll(branchDir, os.ModePerm); err != nil {
		return fmt.Errorf("NewBranch: %w", err)
	}
	if err := os.WriteFile(branchPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("NewBranch: %w", err)
	}
	return nil
}

func ListBranches() ([]Branch, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return nil, fmt.Errorf("ListBranches: %w", err)
	}
	var branches []Branch
	headsDir := filepath.Join(repoDir, "refs", "heads")
	err = filepath.WalkDir(headsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(headsDir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		branches = append(branches, Branch{
			Name:       relPath,
			CommitHash: string(data),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ListBranches: %w", err)
	}
	return branches, nil
}
