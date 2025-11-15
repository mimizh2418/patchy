package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"patchy/objects/objecttype"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"strconv"
	"strings"
)

var objCache = make(map[string][]byte)
var objTypeCache = make(map[string]objecttype.ObjectType)

func objectExists(hash string) bool {
	repoDir, err := repo.FindRepoDir()
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

func WriteObject(objType objecttype.ObjectType, data []byte) (string, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return "", fmt.Errorf("WriteObject: %w", err)
	}
	header := []byte(fmt.Sprintf("%s %d\000", objType.String(), len(data)))
	contents := append(header, data...)
	hash := computeHash(contents)
	if objectExists(hash) {
		objCache[hash] = data
		objTypeCache[hash] = objType
		return hash, nil
	}

	dir := filepath.Join(repoDir, "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])
	compressedData, err := compressObject(contents)
	if err != nil {
		return "", fmt.Errorf("WriteObject: %w", err)
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("WriteObject: %w", err)
	}
	if err := os.WriteFile(file, compressedData, 0666); err != nil {
		return "", fmt.Errorf("WriteObject: %w", err)
	}

	objCache[hash] = data
	objTypeCache[hash] = objType
	return hash, nil
}

func ReadObjectType(hash string) (objecttype.ObjectType, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return objecttype.Unknown, fmt.Errorf("ReadObjectType: %w", err)
	}
	if err := ResolveAndValidateObject(&hash); err != nil {
		return objecttype.Unknown, fmt.Errorf("ReadObjectType: %w", err)
	}

	dir := filepath.Join(repoDir, "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])
	compressedData, err := os.ReadFile(file)
	if err != nil {
		return objecttype.Unknown, fmt.Errorf("ReadObjectType: %w", err)
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return objecttype.Unknown, fmt.Errorf("ReadObjectType: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	headerBytes := make([]byte, 0)
	buf := make([]byte, 1)
	n, err := reader.Read(buf)
	for buf[0] != 0 {
		if n != 1 || err != nil {
			return objecttype.Unknown, fmt.Errorf(
				"ReadObjectType: %w", &BadObject{hash, "format"})
		}
		headerBytes = append(headerBytes, buf[0])
		n, err = reader.Read(buf)
	}
	header := strings.Split(string(headerBytes), " ")
	if len(header) != 2 {
		return objecttype.Unknown, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "header"})
	}

	switch header[0] {
	case "blob":
		return objecttype.Blob, nil
	case "tree":
		return objecttype.Tree, nil
	case "commit":
		return objecttype.Commit, nil
	default:
		return objecttype.Unknown, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "type"})
	}
}

func ReadObject(hash string) (objecttype.ObjectType, []byte, error) {
	repoDir, err := repo.FindRepoDir()
	if err != nil {
		return objecttype.Unknown, nil, fmt.Errorf("ReadObject: %w", err)
	}
	if err := ResolveAndValidateObject(&hash); err != nil {
		return objecttype.Unknown, nil, fmt.Errorf("ReadObject: %w", err)
	}
	if data, ok := objCache[hash]; ok {
		return objTypeCache[hash], data, nil
	}

	dir := filepath.Join(repoDir, "objects", hash[:2])
	file := filepath.Join(dir, hash[2:])
	compressedData, err := os.ReadFile(file)
	if err != nil {
		return objecttype.Unknown, nil, fmt.Errorf("ReadObject: %w", err)
	}

	blob, err := decompressObject(compressedData)
	if err != nil {
		return objecttype.Unknown, nil, fmt.Errorf("ReadObject: %w", err)
	}

	nullPos := -1
	for i, b := range blob {
		if b == 0 {
			nullPos = i
			break
		}
	}
	if nullPos <= 0 {
		return objecttype.Unknown, nil, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "format"})
	}
	header := strings.Split(string(blob[:nullPos]), " ")
	content := blob[nullPos+1:]

	if len(header) != 2 {
		return objecttype.Unknown, nil, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "header"})
	}
	length, err := strconv.Atoi(header[1])
	if err != nil || length != len(content) {
		return objecttype.Unknown, nil, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "header"})
	}

	var objType objecttype.ObjectType
	switch header[0] {
	case "blob":
		objType = objecttype.Blob
	case "tree":
		objType = objecttype.Tree
	case "commit":
		objType = objecttype.Commit
	default:
		return objecttype.Unknown, nil, fmt.Errorf(
			"ReadObjectType: %w", &BadObject{hash, "type"})
	}
	objCache[hash] = content
	objTypeCache[hash] = objType
	return objType, content, nil
}
