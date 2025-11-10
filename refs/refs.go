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
		return err
	}
	if objType, err := objects.ReadObjectType(commitHash); err == nil && objType != objecttype.Commit {
		return fmt.Errorf("object %s is not a commit", commitHash)
	} else if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(repoDir, ref), []byte(commitHash), 0666); err != nil {
		return err
	}
	return nil
}

func ResolveRef(ref string) (string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return "", err
	}
	if data, err := os.ReadFile(filepath.Join(repoDir, ref)); err == nil {
		hash := string(data)
		if objType, err := objects.ReadObjectType(hash); err == nil && objType != objecttype.Commit {
			return "", fmt.Errorf("object %s is not a commit", hash)
		} else if err != nil {
			return "", err
		}
		return hash, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return "", nil
	} else {
		return "", err
	}
}

func ReadHead() (bool, string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return false, "", err
	}
	data, err := os.ReadFile(filepath.Join(repoDir, "HEAD"))
	if err != nil {
		return false, "", err
	}
	content := strings.Split(string(data), "\n")[0]
	if strings.HasPrefix(content, "ref: ") {
		return false, strings.TrimPrefix(content, "ref: "), nil
	}

	if objType, err := objects.ReadObjectType(content); err == nil && objType != objecttype.Commit {
		return false, "", fmt.Errorf("object %s is not a commit", content)
	} else if err != nil {
		return false, "", err
	}
	return true, content, nil
}
