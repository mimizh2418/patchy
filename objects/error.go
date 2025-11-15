package objects

import (
	"patchy/objects/objecttype"
	"strings"
)

type ObjectNotFound struct {
	Hash string
}

func (e *ObjectNotFound) Error() string {
	return "object not found: " + e.Hash
}

type BadObjectID struct {
	Hash string
}

func (e *BadObjectID) Error() string {
	return e.Hash + " is not a valid object id"
}

type AmbiguousObjectID struct {
	ShortHash string
	Hashes    []string
}

func (e *AmbiguousObjectID) Error() string {
	return "ambiguous object id " + e.ShortHash + "; could refer to:\n  " + strings.Join(e.Hashes, "\n  ")
}

type ObjectTypeMismatch struct {
	Hash     string
	Expected objecttype.ObjectType
	Actual   objecttype.ObjectType
}

func (e *ObjectTypeMismatch) Error() string {
	return "object " + e.Hash + " is of type " + e.Actual.String() + ", expected " + e.Expected.String()
}

type BadObject struct {
	Hash        string
	Description string
}

func (e *BadObject) Error() string {
	return e.Hash + " has invalid " + e.Description
}

var (
	ErrObjectNotFound     *ObjectNotFound
	ErrBadObjectID        *BadObjectID
	ErrAmbiguousObjectID  *AmbiguousObjectID
	ErrObjectTypeMismatch *ObjectTypeMismatch
	ErrBadObject          *BadObject
)
