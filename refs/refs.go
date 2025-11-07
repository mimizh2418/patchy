package refs

import (
	"fmt"
	"os"
	"patchy/objects"
	"patchy/objects/objecttype"
	"patchy/repo"
	"path/filepath"
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
