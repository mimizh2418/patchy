package refs

import (
	"errors"
	"fmt"
	"os"
	"patchy/objects"
	"patchy/objects/objecttype"
	"patchy/repo"
	"path/filepath"
	"strings"
)

func UpdateRef(ref string, commitHash string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}
	if objType, err := objects.ReadObjectType(commitHash); err == nil && objType != objecttype.Commit {
		return fmt.Errorf(
			"UpdateRef: %w",
			&objects.ErrObjectTypeMismatch{Hash: commitHash, Expected: objecttype.Commit, Actual: objType})
	} else if err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}

	if err := os.WriteFile(filepath.Join(repoDir, ref), []byte(commitHash), 0666); err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}
	return nil
}

func ResolveRef(ref string) (string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return "", fmt.Errorf("UpdateRef: %w", err)
	}
	if data, err := os.ReadFile(filepath.Join(repoDir, ref)); err == nil {
		hash := string(data)
		if objType, err := objects.ReadObjectType(hash); err == nil && objType != objecttype.Commit {
			return "", fmt.Errorf(
				"ResolveRef: %w",
				&objects.ErrObjectTypeMismatch{Hash: hash, Expected: objecttype.Commit, Actual: objType})
		} else if err != nil {
			return "", fmt.Errorf("ResolveRef: %w", err)
		}
		return hash, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return "", nil
	} else {
		return "", fmt.Errorf("ResolveRef: %w", err)
	}
}

func ReadHead() (bool, string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return false, "", fmt.Errorf("UpdateRef: %w", err)
	}
	data, err := os.ReadFile(filepath.Join(repoDir, "HEAD"))
	if err != nil {
		return false, "", fmt.Errorf("ReadHead: %w", err)

	}
	content := strings.Split(string(data), "\n")[0]
	if strings.HasPrefix(content, "ref: ") {
		return false, strings.TrimPrefix(content, "ref: "), nil
	}

	if objType, err := objects.ReadObjectType(content); err == nil && objType != objecttype.Commit {
		return false, "", fmt.Errorf(
			"ReadHead: %w",
			&objects.ErrObjectTypeMismatch{Hash: content, Expected: objecttype.Commit, Actual: objType})
	} else if err != nil {
		return false, "", fmt.Errorf("ReadHead: %w", err)
	}
	return true, content, nil
}
