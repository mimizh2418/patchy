package refs

import (
	"errors"
	"fmt"
	"os"
	"patchy/objects"
	"patchy/objects/objecttype"
	"patchy/repo"
	"path/filepath"
	"strconv"
	"strings"
)

type HeadState struct {
	Detached bool
	Ref      string
	Commit   string
}

func ResolveRef(ref string) (string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return "", fmt.Errorf("ResolveRef: %w", err)
	}
	if data, err := os.ReadFile(filepath.Join(repoDir, ref)); err == nil {
		hash := string(data)
		if objType, err := objects.ReadObjectType(hash); err == nil && objType != objecttype.Commit {
			return "", fmt.Errorf(
				"ResolveRef: %w",
				&objects.ObjectTypeMismatch{Hash: hash, Expected: objecttype.Commit, Actual: objType})
		} else if err != nil {
			return "", fmt.Errorf("ResolveRef: %w", err)
		}
		return hash, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("ResolveRef: %w", &InvalidRef{Ref: ref})
	} else {
		return "", fmt.Errorf("ResolveRef: %w", err)
	}
}

func ParseRev(revSpec string) (string, error) { // TODO make better name
	commit, err := ResolveRef("refs/heads/" + revSpec)
	if err == nil {
		return commit, nil
	}
	prefix := ""
	if strings.HasPrefix(revSpec, "HEAD") {
		prefix = "HEAD"
	} else if strings.HasPrefix(revSpec, "@") {
		prefix = "@"
	}
	if prefix != "" {
		suffix := strings.TrimPrefix(revSpec, prefix)
		head, err := ReadHead()
		if err != nil {
			return "", fmt.Errorf("ParseRev: %w", err)
		}
		if len(suffix) == 0 {
			return head.Commit, nil
		} else if !strings.HasPrefix(suffix, "~") || strings.HasPrefix(suffix, "^") {
			return "", fmt.Errorf("ParseRev: %w", &InvalidRevSpec{RevSpec: revSpec})
		}
		numStr := suffix[1:]
		num := 1
		if len(numStr) > 0 {
			num, err = strconv.Atoi(numStr)
			if err != nil || num < 0 {
				return "", fmt.Errorf("ParseRev: %w", &InvalidRevSpec{RevSpec: revSpec})
			}
		}
		currentHash := head.Commit
		for i := 0; i < num; i++ {
			commitObj, err := objects.ReadCommit(currentHash)
			if err != nil {
				return "", fmt.Errorf("ParseRev: %w", err)
			}
			if commitObj.Parent == nil {
				return "", fmt.Errorf("ParseRev: %w", &InvalidRevSpec{RevSpec: revSpec})
			}
			currentHash = *commitObj.Parent
		}
		return currentHash, nil
	}
	hash := revSpec
	if err := objects.ResolveAndValidateObject(&hash); err == nil {
		if objType, err := objects.ReadObjectType(hash); err == nil && objType == objecttype.Commit {
			return hash, nil
		} else if err != nil {
			return "", fmt.Errorf("ParseRev: %w", err)
		} else {
			return "", fmt.Errorf(
				"ParseRev: %w",
				&objects.ObjectTypeMismatch{Hash: revSpec, Expected: objecttype.Commit, Actual: objType})
		}
	}
	return "", fmt.Errorf("ParseRev: %w", &InvalidRevSpec{RevSpec: revSpec})
}

func UpdateRef(ref string, commitHash string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}
	if objType, err := objects.ReadObjectType(commitHash); err == nil && objType != objecttype.Commit {
		return fmt.Errorf(
			"UpdateRef: %w",
			&objects.ObjectTypeMismatch{Hash: commitHash, Expected: objecttype.Commit, Actual: objType})
	} else if err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}

	if err := os.WriteFile(filepath.Join(repoDir, ref), []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("UpdateRef: %w", err)
	}
	return nil
}

func ReadHead() (*HeadState, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return nil, fmt.Errorf("ReadHead: %w", err)
	}
	data, err := os.ReadFile(filepath.Join(repoDir, "HEAD"))
	if err != nil {
		return nil, fmt.Errorf("ReadHead: %w", err)

	}
	content := strings.Split(string(data), "\n")[0]
	if strings.HasPrefix(content, "ref: ") {
		ref := strings.TrimPrefix(content, "ref: ")
		hash, err := ResolveRef(ref)
		if err != nil && !errors.As(err, &ErrInvalidRef) {
			return nil, fmt.Errorf("ReadHead: %w", err)
		}
		return &HeadState{false, ref, hash}, nil
	}

	if objType, err := objects.ReadObjectType(content); err == nil && objType != objecttype.Commit {
		return nil, fmt.Errorf(
			"ReadHead: %w",
			&objects.ObjectTypeMismatch{Hash: content, Expected: objecttype.Commit, Actual: objType})
	} else if err != nil {
		return nil, fmt.Errorf("ReadHead: %w", err)
	}
	return &HeadState{true, "", content}, nil
}

func UpdateHead(revSpec string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return fmt.Errorf("UpdateHead: %w", err)
	}
	if _, err := ResolveRef("refs/heads/" + revSpec); err == nil {
		// branch
		if err := os.WriteFile(filepath.Join(repoDir, "HEAD"), []byte("ref: refs/heads/"+revSpec), 0666); err != nil {
			return fmt.Errorf("UpdateHead: %w", err)
		}
		return nil
	}
	// commit hash
	hash, err := ParseRev(revSpec)
	if err != nil {
		return fmt.Errorf("UpdateHead: %w", err)
	}
	if err := os.WriteFile(filepath.Join(repoDir, "HEAD"), []byte(hash), 0666); err != nil {
		return fmt.Errorf("UpdateHead: %w", err)
	}
	return nil
}
