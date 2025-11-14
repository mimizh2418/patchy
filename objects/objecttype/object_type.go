package objecttype

type ObjectType int

const (
	Unknown ObjectType = iota
	Blob
	Tree
	Commit
)

func (objType ObjectType) String() string {
	switch objType {
	case Blob:
		return "blob"
	case Tree:
		return "tree"
	case Commit:
		return "commit"
	default:
		return "unknown"
	}
}
