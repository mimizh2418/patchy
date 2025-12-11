package refs

import (
	"fmt"
	"patchy/objects"
	"patchy/repo"
)

func Checkout(revSpec string) error {
	repoRoot, err := repo.FindRepoRoot()
	if err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}

	// Update HEAD to point to the specified revision
	err = UpdateHead(revSpec)
	if err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}

	// Find the commit hash for the specified revision
	commitHash, err := ParseRev(revSpec)
	if err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}
	commit, err := objects.ReadCommit(commitHash)
	if err != nil {
		return fmt.Errorf("Checkout: %w", err)
	}

	// Update the working directory to match the commit's tree
	return objects.UnpackTree(commit.Tree, repoRoot)
}
