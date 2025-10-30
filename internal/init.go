package internal

import (
	"errors"
	"os"
	"patchy/internal/util"
	"path/filepath"
)

func Init(path string) (repoPath string, err error) {
	if _, e := util.FindRepoDir(); e == nil {
		err = errors.New("already inside a repository")
		return
	}

	repoPath = ""

	repoPath, err = filepath.Abs(filepath.Join(path, ".patchy"))
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = os.RemoveAll(repoPath)
		}
	}()

	if err = os.MkdirAll(filepath.Join(repoPath, "objects"), os.ModePerm); err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Join(repoPath, "refs", "heads"), os.ModePerm); err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Join(repoPath, "refs", "tags"), os.ModePerm); err != nil {
		return
	}

	if err = os.WriteFile(filepath.Join(repoPath, "HEAD"), []byte("ref: refs/heads/main\n"), 0666); err != nil {
		return
	}

	return
}
