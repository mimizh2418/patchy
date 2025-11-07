package objecttype

type ObjectType int

const (
	Unknown ObjectType = iota
	Blob
	Tree
	Commit
)
