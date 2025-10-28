package internal

import (
    "errors"
    "os"
    "path/filepath"
)

func Init(path string) (repoPath string, err error) {
    repoPath = ""

    repoPath, err = filepath.Abs(filepath.Join(path, ".patchy"))
    if err != nil {
        return
    }

    if _, e := os.Stat(repoPath); e == nil {
        err = errors.New("repository already initialized")
        return
    } else if !errors.Is(e, os.ErrNotExist) {
        err = e
        return
    }

    if err = os.MkdirAll(filepath.Join(repoPath, "objects"), os.ModePerm); err != nil {
        return
    }
    if err = os.MkdirAll(filepath.Join(repoPath, "refs", "heads"), os.ModePerm); err != nil {
        return
    }
    if err = os.MkdirAll(filepath.Join(repoPath, "refs", "tags"), os.ModePerm); err != nil {
        return
    }

    if err = os.WriteFile(filepath.Join(repoPath, "HEAD"), []byte("ref: refs/heads/main\n"), os.ModePerm); err != nil {
        return
    }

    return
}
