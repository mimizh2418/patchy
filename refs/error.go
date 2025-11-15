package refs

type InvalidRef struct {
	Ref string
}

func (e *InvalidRef) Error() string {
	return "could not resolve ref: " + e.Ref
}

type InvalidRevSpec struct {
	RevSpec string
}

func (e *InvalidRevSpec) Error() string {
	return "unknown revision '" + e.RevSpec + "'"
}

var (
	ErrInvalidRef     *InvalidRef
	ErrInvalidRevSpec *InvalidRevSpec
)
