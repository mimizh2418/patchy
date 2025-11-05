package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
)

func DoesFileExist(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func IsFileInRepo(path string) (bool, error) {
	repoRoot, err := FindRepoRoot()
	if err != nil {
		return false, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	if _, err := filepath.Rel(repoRoot, absPath); err != nil {
		return false, nil
	}
	return true, nil
}

func IsDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	} else {
		return info.IsDir(), nil
	}
}

func ReadFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
