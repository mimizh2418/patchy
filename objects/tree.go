package objects

import (
	"encoding/hex"
	"fmt"
	"os"
	"patchy/ignore"
	"patchy/objects/objecttype"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"text/tabwriter"

	"github.com/fatih/color"
)

type TreeEntry struct {
	Mode     string
	Name     string
	Hash     string
	Children []TreeEntry
}

func WriteTree(path string) (string, error) {
	repoRoot, err := repo.FindRepoRoot()
	if err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}

	// Validate path
	if err = repo.ValidateFileInRepo(path); err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}
	if isDir, err := util.IsDirectory(path); err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	} else if !isDir {
		return "", fmt.Errorf("WriteTree: file %s is not a directory", path)
	}

	ignoreList, err := ignore.ReadIgnoreFile()
	if err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}
	entries := make([]TreeEntry, 0)
	err = filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file == path {
			return nil
		}

		absPath, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(repoRoot, absPath)
		if err != nil {
			return err
		}
		if info.IsDir() {
			relPath += "/"
		}
		for _, pattern := range ignoreList {
			matched, err := filepath.Match(pattern, relPath)
			if err != nil {
				continue
			}
			if matched {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		name := filepath.Base(file)
		if info.IsDir() {
			hash, err := WriteTree(file)
			if err != nil {
				return err
			}
			entries = append(entries, TreeEntry{"040000", name, hash, []TreeEntry{}})
			return filepath.SkipDir
		}
		hash, err := WriteBlob(file)
		entries = append(entries, TreeEntry{"100644", name, hash, []TreeEntry{}})
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}
	data := make([]byte, 0)
	for _, entry := range entries {
		entryData := []byte(fmt.Sprintf("%s\000%s\000", entry.Mode, entry.Name))
		rawHash, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return "", fmt.Errorf("WriteTree: %w", err)
		}
		entryData = append(entryData, rawHash...)
		data = append(data, entryData...)
	}

	hash, err := WriteObject(objecttype.Tree, data)
	if err != nil {
		return "", fmt.Errorf("WriteTree: %w", err)
	}
	return hash, nil
}

func ReadTree(hash string) ([]TreeEntry, error) {
	objType, data, err := ReadObject(hash)
	if err != nil {
		return nil, fmt.Errorf("WriteTree: %w", err)
	}
	if objType != objecttype.Tree {
		return nil, fmt.Errorf(
			"WriteTree: %w", &ObjectTypeMismatch{hash, objecttype.Tree, objType})
	}
	entries := make([]TreeEntry, 0)
	i := 0
	for i < len(data) {
		modeEnd := i
		for data[modeEnd] != 0 {
			modeEnd++
			if modeEnd >= len(data) {
				return nil, fmt.Errorf("WriteTree: %w", &BadObject{hash, "format"})
			}
		}
		mode := string(data[i:modeEnd])
		i = modeEnd + 1

		nameEnd := i
		for data[nameEnd] != 0 {
			nameEnd++
			if nameEnd >= len(data) {
				return nil, fmt.Errorf("WriteTree: %w", &BadObject{hash, "format"})
			}
		}
		name := string(data[i:nameEnd])
		i = nameEnd + 1

		if i+20 > len(data) {
			return nil, fmt.Errorf("WriteTree: %w", &BadObject{hash, "format"})
		}
		rawHash := data[i : i+20]
		hash := hex.EncodeToString(rawHash)
		i += 20

		entries = append(entries, TreeEntry{mode, name, hash, []TreeEntry{}})
	}
	return entries, nil
}

func ReadTreeRecursive(hash string) ([]TreeEntry, error) {
	entries, err := ReadTree(hash)
	if err != nil {
		return nil, err
	}
	for i, entry := range entries {
		if entry.Mode == "040000" {
			entries[i].Children, err = ReadTreeRecursive(entry.Hash)
			if err != nil {
				return nil, err
			}
		}
	}
	return entries, nil
}

func FlattenTreeEntries(entries []TreeEntry) []TreeEntry {
	flatEntries := make([]TreeEntry, 0)
	for _, entry := range entries {
		if entry.Mode == "040000" {
			children := FlattenTreeEntries(entry.Children)
			for _, child := range children {
				child.Name = filepath.Join(entry.Name, child.Name)
				flatEntries = append(flatEntries, child)
			}
		} else {
			flatEntries = append(flatEntries, entry)
		}
	}
	return flatEntries
}

func PrintTree(hash string) error {
	entries, err := ReadTree(hash)
	if err != nil {
		return fmt.Errorf("PrintTree: %w", err)
	}

	util.ColorPrintf(color.FgCyan, "[tree %s]\n", resolveObject(hash))
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, entry := range entries {
		util.Fprintf(writer, "%s\t%s  \t%s\n", entry.Mode, entry.Hash, entry.Name)
	}
	return writer.Flush()
}

func UnpackTree(hash string, path string) error {
	// Validate path
	if err := repo.ValidateFileInRepo(path); err != nil {
		return fmt.Errorf("UnpackTree: %w", err)
	}
	if isDir, err := util.IsDirectory(path); err != nil {
		return fmt.Errorf("UnpackTree: %w", err)
	} else if !isDir {
		return fmt.Errorf("UnpackTree: file %s is not a directory", path)
	}

	tree, err := ReadTreeRecursive(hash)
	if err != nil {
		return fmt.Errorf("UnpackTree: %w", err)
	}
	entries := FlattenTreeEntries(tree)
	for _, entry := range entries {
		file := filepath.Join(path, entry.Name)
		blob, err := ReadBlob(entry.Hash)
		if err != nil {
			return fmt.Errorf("UnpackTree: %w", err)
		}
		if err = os.WriteFile(file, blob, 0644); err != nil {
			return fmt.Errorf("UnpackTree: %w", err)
		}
	}
	return nil
}
