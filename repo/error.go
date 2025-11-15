package repo

import "errors"

type FileNotInRepo struct {
	Path string
}

func (e *FileNotInRepo) Error() string {
	return "file " + e.Path + " is not inside this repository"
}

var (
	ErrAlreadyInRepo = errors.New("current directory is already part of a repository")
	ErrNotInRepo     = errors.New("current directory is not inside of a repository")
	ErrFileNotInRepo *FileNotInRepo
)
