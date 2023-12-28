package errors

import "errors"

var (
	New = errors.New
	Is  = errors.Is
	Ad  = errors.As
)

func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	Err error
	Msg string
}

func (p publicError) Error() string {
	return p.Err.Error()
}

func (p publicError) Public() string {
	return p.Msg
}
func (p publicError) Unwrap() error {
	return p.Err
}
