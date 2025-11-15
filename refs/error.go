package refs

type InvalidRef struct {
	Ref string
}

func (e *InvalidRef) Error() string {
	return "could not resolve ref: " + e.Ref
}

var ErrInvalidRef *InvalidRef
