package objects

import (
	"patchy/objects/objecttype"
	"strings"
)

type ErrObjectNotFound struct {
	Hash string
}

func (e *ErrObjectNotFound) Error() string {
	return "object not found: " + e.Hash
}

type ErrBadObjectID struct {
	Hash string
}

func (e *ErrBadObjectID) Error() string {
	return e.Hash + " is not a valid object id"
}

type ErrAmbiguousObjectID struct {
	ShortHash string
	Hashes    []string
}

func (e *ErrAmbiguousObjectID) Error() string {
	return "ambiguous object id " + e.ShortHash + "; could refer to:\n  " + strings.Join(e.Hashes, "\n  ")
}

type ErrObjectTypeMismatch struct {
	Hash     string
	Expected objecttype.ObjectType
	Actual   objecttype.ObjectType
}

func (e *ErrObjectTypeMismatch) Error() string {
	return "object " + e.Hash + " is of type " + e.Actual.String() + ", expected " + e.Expected.String()
}

type ErrBadObject struct {
	Hash        string
	Description string
}

func (e *ErrBadObject) Error() string {
	return e.Hash + " has invalid " + e.Description
}
