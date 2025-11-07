package ignore

import (
	"bufio"
	"os"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"strings"
)

var ignorePatterns []string = nil

func ReadIgnoreFile() ([]string, error) {
	if ignorePatterns != nil {
		return ignorePatterns, nil
	}

	repoRoot, err := repo.FindRepoRoot()
	if err != nil {
		return nil, err
	}

	patterns := []string{".patchy/", ".git/"}
	ignoreFileExists, err := util.DoesFileExist(filepath.Join(repoRoot, ".patchyignore"))
	if err != nil {
		return nil, err
	}
	if !ignoreFileExists {
		return patterns, nil
	}

	f, err := os.Open(filepath.Join(repoRoot, ".patchyignore"))
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	ignorePatterns = patterns
	return patterns, nil
}
