package objects

import (
    "encoding/hex"
    "errors"
    "fmt"
    "os"
    "patchy/ignore"
    "patchy/objects/objecttype"
    "patchy/repo"
    "patchy/util"
    "path/filepath"
    "text/tabwriter"
)

type TreeEntry struct {
    Mode string
    Name string
    Hash string
}

func WriteTree(path string) (string, error) {
    repoDir, err := repo.FindRepoDir()
    if err != nil {
        return "", err
    }
    repoRoot := filepath.Dir(repoDir)

    // Validate path
    if exists, err := util.DoesFileExist(path); err != nil {
        return "", err
    } else if !exists {
        return "", nil
    }
    if inRepo, err := repo.IsFileInRepo(path); err != nil {
        return "", err
    } else if !inRepo {
        return "", errors.New("file not in this repository")
    }
    if isDir, err := util.IsDirectory(path); err != nil {
        return "", err
    } else if !isDir {
        return "", errors.New("not a directory")
    }

    ignoreList, err := ignore.ReadIgnoreFile()
    if err != nil {
        return "", err
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
            entries = append(entries, TreeEntry{"040000", name, hash})
            return filepath.SkipDir
        }
        hash, data, err := HashObject(file)
        if err != nil {
            return err
        }
        if !objectExists(hash) {
            if err = WriteObject(hash, data); err != nil {
                return err
            }
        }
        entries = append(entries, TreeEntry{"100644", name, hash})
        return nil
    })
    if err != nil {
        return "", err
    }
    data := make([]byte, 0)
    for _, entry := range entries {
        entryData := []byte(fmt.Sprintf("%s\000%s\000", entry.Mode, entry.Name))
        rawHash, err := hex.DecodeString(entry.Hash)
        if err != nil {
            return "", err
        }
        entryData = append(entryData, rawHash...)
        data = append(data, entryData...)
    }
    header := []byte(fmt.Sprintf("tree %d\000", len(data)))
    data = append(header, data...)

    hash := computeHash(data)
    if objectExists(hash) {
        return hash, nil
    }
    if err = WriteObject(hash, data); err != nil {
        return "", err
    }
    return hash, nil
}

func ReadTree(hash string) ([]TreeEntry, error) {
    objType, data, err := ReadObject(hash)
    if err != nil {
        return nil, err
    }
    if objType != objecttype.Tree {
        return nil, fmt.Errorf("object %s is not a tree", hash)
    }
    entries := make([]TreeEntry, 0)
    i := 0
    for i < len(data) {
        modeEnd := i
        for data[modeEnd] != 0 {
            modeEnd++
            if modeEnd >= len(data) {
                return nil, errors.New("invalid tree format")
            }
        }
        mode := string(data[i:modeEnd])
        i = modeEnd + 1

        nameEnd := i
        for data[nameEnd] != 0 {
            nameEnd++
            if nameEnd >= len(data) {
                return nil, errors.New("invalid tree format")
            }
        }
        name := string(data[i:nameEnd])
        i = nameEnd + 1

        if i+20 > len(data) {
            return nil, errors.New("invalid tree format")
        }
        rawHash := data[i : i+20]
        hash := hex.EncodeToString(rawHash)
        i += 20

        entries = append(entries, TreeEntry{mode, name, hash})
    }
    return entries, nil
}

func PrintTree(hash string) error {
    entries, err := ReadTree(hash)
    if err != nil {
        return err
    }

    writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    for _, entry := range entries {
        util.Fprintf(writer, "%s\t%s  \t%s\n", entry.Mode, entry.Hash, entry.Name)
    }
    return writer.Flush()
}
