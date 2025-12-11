package diff

import (
	"fmt"
	"patchy/objects"
	"patchy/refs"
	"patchy/repo"
	"patchy/util"
	"sort"

	"github.com/fatih/color"
)

type ChangeType int

const (
	Added ChangeType = iota
	Deleted
	Modified
	Moved
)

type FileChange struct {
	OldName    string
	NewName    string
	OldHash    string
	NewHash    string
	ChangeType ChangeType
}

func TreeDiff(newTree string, oldTree string) ([]FileChange, error) {
	newEntries, err := objects.ReadTreeRecursive(newTree)
	if err != nil {
		return nil, fmt.Errorf("TreeDiff: %w", err)
	}
	newEntries = objects.FlattenTreeEntries(newEntries)

	var oldEntries []objects.TreeEntry
	if len(oldTree) > 0 {
		oldEntries, err = objects.ReadTreeRecursive(oldTree)
		if err != nil {
			return nil, fmt.Errorf("TreeDiff: %w", err)
		}
		oldEntries = objects.FlattenTreeEntries(oldEntries)
	} else {
		oldEntries = make([]objects.TreeEntry, 0)
	}

	changes := make([]FileChange, 0)
	newTreeByName := make(map[string]string)
	for _, entry := range newEntries {
		newTreeByName[entry.Name] = entry.Hash
	}
	oldTreeByName := make(map[string]string)
	for _, entry := range oldEntries {
		oldTreeByName[entry.Name] = entry.Hash
	}
	newTreeByHash := make(map[string]string)
	for _, entry := range newEntries {
		newTreeByHash[entry.Hash] = entry.Name
	}
	oldTreeByHash := make(map[string]string)
	for _, entry := range oldEntries {
		oldTreeByHash[entry.Hash] = entry.Name
	}

	for name, hash1 := range newTreeByName {
		if hash2, exists := oldTreeByName[name]; exists {
			if hash1 != hash2 {
				changes = append(changes, FileChange{
					OldName:    name,
					NewName:    name,
					OldHash:    hash1,
					NewHash:    hash2,
					ChangeType: Modified,
				})
			}
		} else if _, renamed := oldTreeByHash[hash1]; !renamed {
			changes = append(changes, FileChange{
				OldName:    "",
				NewName:    name,
				OldHash:    "",
				NewHash:    hash1,
				ChangeType: Added,
			})
		}
	}
	// Check for deleted files
	for name, hash2 := range oldTreeByName {
		_, exists := newTreeByName[name]
		_, renamed := newTreeByHash[hash2]
		if !exists && !renamed {
			changes = append(changes, FileChange{
				OldName:    name,
				NewName:    "",
				OldHash:    hash2,
				NewHash:    "",
				ChangeType: Deleted,
			})
		}
	}
	// Check for renamed/moved files
	for hash, oldName := range oldTreeByHash {
		if newName, exists := newTreeByHash[hash]; exists && oldName != newName {
			changes = append(changes, FileChange{
				OldName:    oldName,
				NewName:    newName,
				OldHash:    hash,
				NewHash:    hash,
				ChangeType: Moved,
			})
		}
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].NewName < changes[j].NewName
	})
	return changes, nil
}

func WorkingTreeDiff() ([]FileChange, error) {
	repoRoot, err := repo.FindRepoRoot()
	if err != nil {
		return nil, fmt.Errorf("WorkingTreeDiff: %w", err)
	}
	headState, err := refs.ReadHead()
	if err != nil {
		return nil, fmt.Errorf("WorkingTreeDiff: %w", err)
	}
	headCommit, err := objects.ReadCommit(headState.Commit)
	if err != nil {
		return nil, fmt.Errorf("WorkingTreeDiff: %w", err)
	}
	currentTree, err := objects.WriteTree(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("WorkingTreeDiff: %w", err)
	}
	changes, err := TreeDiff(currentTree, headCommit.Tree)
	if err != nil {
		return nil, fmt.Errorf("WorkingTreeDiff: %w", err)
	}
	return changes, nil
}

func PrintDiffSummary(changes []FileChange) {
	additions := 0
	modifications := 0
	deletions := 0
	moves := 0
	for _, change := range changes {
		switch change.ChangeType {
		case Added:
			additions++
		case Deleted:
			deletions++
		case Modified:
			modifications++
		case Moved:
			moves++
		}
	}
	if additions > 0 {
		util.ColorPrintf(color.FgGreen, "    %d new file(s)\n", additions)
	}
	if modifications > 0 {
		util.ColorPrintf(color.FgYellow, "    %d file(s) changed\n", modifications)
	}
	if deletions > 0 {
		util.ColorPrintf(color.FgRed, "    %d file(s) removed\n", deletions)
	}
	if moves > 0 {
		util.ColorPrintf(color.FgCyan, "    %d file(s) renamed/moved\n", moves)
	}
}
