package objects

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os/user"
	"patchy/objects/objecttype"
	"patchy/util"
	"strconv"
	"time"
)

type Commit struct {
	Tree    string
	Author  string
	Message string
	Time    time.Time
	Parent  *string
}

func WriteCommit(tree string, parent *string, message string) (string, error) {
	if !objectExists(tree) {
		return "", fmt.Errorf("tree object %s does not exist", tree)
	}
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	author := currentUser.Username
	currentTime := time.Now()
	data, err := hex.DecodeString(tree)
	if err != nil {
		return "", err
	}
	data = append(data, []byte(fmt.Sprintf("\000%s\000%s\000%d\000", author, message, currentTime.Unix()))...)
	if parent != nil {
		if objType, err := ReadObjectType(*parent); err == nil && objType != objecttype.Commit {
			return "", fmt.Errorf("invalid parent, object %s is not a commit", *parent)
		} else if err != nil {
			return "", err
		}
		rawParentHash, err := hex.DecodeString(*parent)
		if err != nil {
			return "", err
		}
		data = append(data, rawParentHash...)
	}
	header := []byte(fmt.Sprintf("commit %d\000", len(data)))
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

func ReadCommit(hash string) (*Commit, error) {
	objType, data, err := ReadObject(hash)
	if err != nil {
		return nil, err
	}
	if objType != objecttype.Commit {
		return nil, fmt.Errorf("object %s is not a commit", hash)
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
		return nil, errors.New("invalid commit format")
	}
	treeHash := hex.EncodeToString(data[:treeHashEnd])
	if len(treeHash) != 40 || !objectExists(treeHash) {
		return nil, fmt.Errorf("invalid tree hash %s in commit", treeHash)
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
		return nil, fmt.Errorf("invalid commit format")
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
		return nil, fmt.Errorf("invalid commit format")
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
		return nil, fmt.Errorf("invalid commit format")
	}
	unixTime, err := strconv.Atoi(string(data[messageEnd+1 : timeEnd]))
	if err != nil {
		return nil, errors.New("invalid commit format")
	}
	commit.Time = time.Unix(int64(unixTime), 0)
	i++

	commit.Parent = nil
	if i < len(data) {
		parentHash := hex.EncodeToString(data[timeEnd+1:])
		if len(parentHash) != 40 {
			return nil, errors.New("invalid commit format")
		}
		if objType, err := ReadObjectType(parentHash); err == nil && objType != objecttype.Commit {
			return nil, fmt.Errorf("invalid parent, object %s is not a commit", parentHash)
		} else if err != nil {
			return nil, err
		}
		commit.Parent = &parentHash
	}
	return commit, nil
}

func PrintCommit(hash string) error {
	commit, err := ReadCommit(hash)
	if err != nil {
		return err
	}
	util.Printf("commit %s\n", hash)
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
