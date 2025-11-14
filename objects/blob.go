package objects

import (
	"fmt"
	"os"
	"patchy/objects/objecttype"
	"patchy/util"
)

func WriteBlob(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	hash, err := WriteObject(objecttype.Blob, data)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func ReadBlob(hash string) ([]byte, error) {
	objType, blob, err := ReadObject(hash)
	if err != nil {
		return nil, err
	}
	if objType != objecttype.Blob {
		return nil, fmt.Errorf("object %s is not a blob", hash)
	}
	return blob, nil
}

func PrintBlob(hash string) error {
	data, err := ReadBlob(hash)
	if err != nil {
		return err
	}
	util.Println(string(data))
	return nil
}
