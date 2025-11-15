package refs

type ErrInvalidRef struct {
	Ref string
}

func (e *ErrInvalidRef) Error() string {
	return "could not resolve ref: " + e.Ref
}
