package objects

import (
    "bytes"
    "compress/zlib"
    "crypto/sha1"
    "errors"
    "fmt"
    "io"
    "os"
    "patchy/internal/util"
    "path/filepath"
    "strconv"
    "strings"
)

type ObjectType int

const (
    Unknown ObjectType = iota
    Blob
    Tree
    Commit
)

func doesObjectExist(hash string) bool {
    repoDir, err := util.FindRepoDir()
    if err != nil {
        return false
    }
    exists, err := util.DoesFileExist(filepath.Join(repoDir, "objects", hash[:2], hash[2:]))
    return exists && err == nil
}

func compressObject(object []byte) ([]byte, error) {
    var data bytes.Buffer
    writer := zlib.NewWriter(&data)
    if _, err := writer.Write(object); err != nil {
        _ = writer.Close()
        return nil, err
    }
    _ = writer.Close()
    return data.Bytes(), nil
}

func decompressObject(compressedObject []byte) ([]byte, error) {
    reader, err := zlib.NewReader(bytes.NewReader(compressedObject))
    if err != nil {
        return nil, err
    }
    defer func() {
        _ = reader.Close()
    }()
    obj, err := io.ReadAll(reader)
    if err != nil {
        return nil, err
    }
    return obj, nil
}

func computeHash(data []byte) string {
    return fmt.Sprintf("%x", sha1.Sum(data))
}

func HashObject(filename string) (string, []byte, error) {
    blob, err := makeBlob(filename)
    if err != nil {
        return "", nil, err
    }
    return computeHash(blob), blob, nil
}

func WriteObject(hash string, data []byte) error {
    repoDir, err := util.FindRepoDir()
    if err != nil {
        return err
    }

    if len(hash) != 40 {
        return errors.New("invalid hash")
    }
    dir := filepath.Join(repoDir, "objects", hash[:2])
    file := filepath.Join(dir, hash[2:])

    compressedData, err := compressObject(data)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return err
    }
    if err := os.WriteFile(file, compressedData, 0666); err != nil {
        return err
    }
    return nil
}

func ReadObjectType(hash string) (ObjectType, error) {
    repoDir, err := util.FindRepoDir()
    if err != nil {
        return Unknown, err
    }

    dir := filepath.Join(repoDir, "objects", hash[:2])
    file := filepath.Join(dir, hash[2:])
    compressedData, err := os.ReadFile(file)
    if errors.Is(err, os.ErrNotExist) {
        return Unknown, fmt.Errorf("object %s not found", hash)
    } else if err != nil {
        return Unknown, err
    }

    reader, err := zlib.NewReader(bytes.NewReader(compressedData))
    if err != nil {
        return Unknown, err
    }
    defer func() {
        _ = reader.Close()
    }()
    headerBytes := make([]byte, 0)
    buf := make([]byte, 1)
    n, err := reader.Read(buf)
    for buf[0] != 0 {
        if n != 1 || err != nil {
            return Unknown, errors.New("invalid object format")
        }
        headerBytes = append(headerBytes, buf[0])
        n, err = reader.Read(buf)
    }
    header := strings.Split(string(headerBytes), " ")
    if len(header) != 2 {
        return Unknown, errors.New("invalid object header")
    }

    switch header[0] {
    case "blob":
        return Blob, nil
    case "tree":
        return Tree, nil
    case "commit":
        return Commit, nil
    default:
        return Unknown, errors.New("invalid object type")
    }
}

func ReadObject(hash string) (ObjectType, []byte, error) {
    repoDir, err := util.FindRepoDir()
    if err != nil {
        return Unknown, nil, err
    }

    dir := filepath.Join(repoDir, "objects", hash[:2])
    file := filepath.Join(dir, hash[2:])
    compressedData, err := os.ReadFile(file)
    if errors.Is(err, os.ErrNotExist) {
        return Unknown, nil, fmt.Errorf("object %s not found", hash)
    } else if err != nil {
        return Unknown, nil, err
    }

    blob, err := decompressObject(compressedData)
    if err != nil {
        return Unknown, nil, err
    }

    nullPos := -1
    for i, b := range blob {
        if b == 0 {
            nullPos = i
            break
        }
    }
    if nullPos <= 0 {
        return Unknown, nil, errors.New("invalid object format")
    }
    header := strings.Split(string(blob[:nullPos]), " ")
    content := blob[nullPos+1:]

    if len(header) != 2 {
        return Unknown, nil, errors.New("invalid object header")
    }
    length, err := strconv.Atoi(header[1])
    if err != nil {
        return Unknown, nil, errors.New("invalid object header")
    } else if length != len(content) {
        return Unknown, nil, errors.New("object length mismatch")
    }

    var objType ObjectType
    switch header[0] {
    case "blob":
        objType = Blob
    case "tree":
        objType = Tree
    case "commit":
        objType = Commit
    default:
        return Unknown, nil, errors.New("invalid object type")
    }
    return objType, content, nil
}
