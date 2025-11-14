package objects

import (
	"encoding/hex"
	"fmt"
	"os/user"
	"patchy/objects/objecttype"
	"patchy/util"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type Commit struct {
	Tree    string
	Author  string
	Message string
	Time    time.Time
	Parent  *string
}

func WriteCommit(tree string, parent *string, message string) (string, error) {
	if err := resolveAndValidateObject(&tree); err != nil {
		return "", fmt.Errorf("WriteCommit: bad tree, %w", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("WriteCommit: %w", err)
	}
	author := currentUser.Username
	currentTime := time.Now()
	data, err := hex.DecodeString(tree)
	if err != nil {
		return "", fmt.Errorf("WriteCommit: %w", err)
	}
	data = append(data, []byte(fmt.Sprintf("\000%s\000%s\000%d\000", author, message, currentTime.Unix()))...)
	if parent != nil {
		if objType, err := ReadObjectType(*parent); err == nil && objType != objecttype.Commit {
			return "", fmt.Errorf("WriteCommit: bad parent, %w ", &ErrObjectTypeMismatch{*parent, objecttype.Commit, objType})
		} else if err != nil {
			return "", fmt.Errorf("WriteCommit: bad parent, %w", err)
		}
		rawParentHash, err := hex.DecodeString(*parent)
		if err != nil {
			return "", fmt.Errorf("WriteCommit: %w", err)
		}
		data = append(data, rawParentHash...)
	}
	hash, err := WriteObject(objecttype.Commit, data)
	if err != nil {
		return "", fmt.Errorf("WriteCommit: %w", err)
	}
	return hash, nil
}

func ReadCommit(hash string) (*Commit, error) {
	objType, data, err := ReadObject(hash)
	if err != nil {
		return nil, fmt.Errorf("ReadCommit: %w", err)
	}
	if objType != objecttype.Commit {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrObjectTypeMismatch{hash, objecttype.Commit, objType})
	}
	i := 0
	treeHashEnd := -1
	commit := &Commit{}
	for ; i < len(data); i++ {
		if data[i] == 0 {
			treeHashEnd = i
			break
		}
	}
	if treeHashEnd == -1 {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrBadObject{hash, "format"})
	}
	treeHash := hex.EncodeToString(data[:treeHashEnd])
	if err := validateObject(treeHash); err != nil {
		return nil, fmt.Errorf("ReadCommit: bad tree, %w", err)
	}
	commit.Tree = treeHash
	i++

	authorEnd := -1
	for ; i < len(data); i++ {
		if data[i] == 0 {
			authorEnd = i
			break
		}
	}
	if authorEnd == -1 {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrBadObject{hash, "format"})
	}
	commit.Author = string(data[treeHashEnd+1 : authorEnd])
	i++

	messageEnd := -1
	for ; i < len(data); i++ {
		if data[i] == 0 {
			messageEnd = i
			break
		}
	}
	if messageEnd == -1 {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrBadObject{hash, "format"})
	}
	commit.Message = string(data[authorEnd+1 : messageEnd])
	i++

	timeEnd := -1
	for ; i < len(data); i++ {
		if data[i] == 0 {
			timeEnd = i
			break
		}
	}
	if timeEnd == -1 {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrBadObject{hash, "format"})
	}
	unixTime, err := strconv.Atoi(string(data[messageEnd+1 : timeEnd]))
	if err != nil {
		return nil, fmt.Errorf("ReadCommit: %w", &ErrBadObject{hash, "format"})
	}
	commit.Time = time.Unix(int64(unixTime), 0)
	i++

	commit.Parent = nil
	if i < len(data) {
		parentHash := hex.EncodeToString(data[timeEnd+1:])
		if len(parentHash) != 40 {
			return nil, fmt.Errorf("ReadCommit: bad parent, %w", &ErrBadObjectID{hash})
		}
		if objType, err := ReadObjectType(parentHash); err == nil && objType != objecttype.Commit {
			return nil, fmt.Errorf("ReadCommit: bad parent, %w ", &ErrObjectTypeMismatch{parentHash, objecttype.Commit, objType})
		} else if err != nil {
			return nil, fmt.Errorf("ReadCommit: %w", err)
		}
		commit.Parent = &parentHash
	}
	return commit, nil
}

func PrintCommit(hash string) error {
	commit, err := ReadCommit(hash)
	if err != nil {
		return fmt.Errorf("PrintCommit: %w", err)
	}
	util.ColorPrintf(color.FgCyan, "[commit %s]\n", resolveObject(hash))
	util.Printf("tree %s\n", commit.Tree)
	if commit.Parent != nil {
		util.Printf("parent %s\n", *commit.Parent)
	}
	util.Printf("author %s\n", commit.Author)
	util.Printf("date %s\n\n", commit.Time.Format(time.RubyDate))
	if commit.Message != "" {
		util.Printf("    %s\n", commit.Message)
	}
	return nil
}
