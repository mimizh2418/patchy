package objects

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"strings"
)

func validateObject(hash string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return err
	}
	if _, err := hex.DecodeString(hash); err != nil || len(hash) != 40 {
		return fmt.Errorf("'%s' is not a valid object id", hash)
	}
	if exists, err := util.DoesFileExist(filepath.Join(repoDir, "objects", hash[:2], hash[2:])); err == nil && !exists {
		return fmt.Errorf("object %s not found", hash)
	} else if err != nil {
		return err
	}
	return nil
}

func resolveObject(shortHash string) string {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return ""
	}
	if exists, err := util.DoesFileExist(filepath.Join(repoDir, "objects", shortHash[:2])); err == nil && !exists {
		return ""
	}
	matches := make([]string, 0)
	objectsDir := filepath.Join(repoDir, "objects", shortHash[:2])
	err = filepath.WalkDir(objectsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), shortHash[2:]) {
			matches = append(matches, shortHash[:2]+d.Name())
		}
		return nil
	})
	if err != nil {
		return ""
	}

	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}

func resolveAndValidateObject(shortHash *string) error {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return err
	}
	decodeCheckString := *shortHash
	if len(decodeCheckString)%2 != 0 {
		decodeCheckString += "0"
	}
	if _, err := hex.DecodeString(decodeCheckString); err != nil || len(*shortHash) < 4 || len(*shortHash) > 40 {
		return fmt.Errorf("'%s' is not a valid object id", *shortHash)
	}
	if exists, err := util.DoesFileExist(filepath.Join(repoDir, "objects", (*shortHash)[:2])); err == nil && !exists {
		return fmt.Errorf("object %s not found", *shortHash)
	}
	matches := make([]string, 0)
	objectsDir := filepath.Join(repoDir, "objects", (*shortHash)[:2])
	err = filepath.WalkDir(objectsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), (*shortHash)[2:]) {
			matches = append(matches, (*shortHash)[:2]+d.Name())
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return fmt.Errorf("object %s not found", *shortHash)
	}
	if len(matches) == 1 {
		*shortHash = matches[0]
		return nil
	}
	return fmt.Errorf("multiple objects found with prefix %s: %v", *shortHash, matches)
}
