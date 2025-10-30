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

func makeBlob(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	header := []byte(fmt.Sprintf("blob %d\000", len(data)))
	return append(header, data...), nil
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

func computeHash(blob []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(blob))
}

func HashObject(filename string) (string, []byte, error) {
	blob, err := makeBlob(filename)
	if err != nil {
		return "", []byte{}, err
	}
	return computeHash(blob), blob, nil
}

func WriteObject(hash string, blob []byte) error {
	repoDir, err := util.FindRepoDir()
	if err != nil {
		return err
	}

	if len(hash) != 40 {
		return errors.New("invalid hash")
	}
	dir := filepath.Join(repoDir, "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])

	compressedBlob, err := compressObject(blob)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(file, compressedBlob, 0666); err != nil {
		return err
	}
	return nil
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
		return Unknown, nil, errors.New("object " + hash + " not found")
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
