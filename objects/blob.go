package objects

import (
	"fmt"
	"os"
	"patchy/objects/objecttype"
	"patchy/util"
)

func makeBlob(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	header := []byte(fmt.Sprintf("blob %d\000", len(data)))
	return append(header, data...), nil
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
