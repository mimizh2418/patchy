package objects

import (
	"fmt"
	"os"
	"patchy/objects/objecttype"
	"patchy/util"

	"github.com/fatih/color"
)

func WriteBlob(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("WriteBlob: %w", err)
	}
	hash, err := WriteObject(objecttype.Blob, data)
	if err != nil {
		return "", fmt.Errorf("WriteBlob: %w", err)
	}
	return hash, nil
}

func ReadBlob(hash string) ([]byte, error) {
	objType, blob, err := ReadObject(hash)
	if err != nil {
		return nil, fmt.Errorf("ReadBlob: %w", err)
	}
	if objType != objecttype.Blob {
		return nil, fmt.Errorf(
			"ReadBlob: %w", &ErrObjectTypeMismatch{hash, objecttype.Blob, objType})
	}
	return blob, nil
}

func PrintBlob(hash string) error {
	data, err := ReadBlob(hash)
	if err != nil {
		return fmt.Errorf("PrintBlob: %w", err)
	}
	util.ColorPrintf(color.FgCyan, "[blob %s]\n", resolveObject(hash))
	util.Println(string(data))
	return nil
}
