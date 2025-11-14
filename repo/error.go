package repo

type ErrAlreadyInRepo struct{}

func (e *ErrAlreadyInRepo) Error() string {
	return "current directory is already part of a repository"
}

type ErrNotInRepo struct{}

func (*ErrNotInRepo) Error() string {
	return "current directory is not inside of a repository"
}

type ErrFileNotInRepo struct {
	Path string
}

func (e *ErrFileNotInRepo) Error() string {
	return "file " + e.Path + " is not inside this repository"
}
