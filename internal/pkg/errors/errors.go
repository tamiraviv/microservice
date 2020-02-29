package errors

import (
	"fmt"
	"io"
)

type ErrorType int

const (
	ErrorTypeUnkown ErrorType = iota
	ErrorTypeBadRequest
	ErrorTypeNotFound
	ErrorTypeInternal
)

type Err struct {
	err       error
	errorType ErrorType
	msg       string
}

func New(s string) *Err {
	return &Err{
		msg: s,
	}
}

func Errorf(s string, a ...interface{}) *Err {
	return &Err{
		msg: fmt.Sprintf(s, a),
	}
}

func (e *Err) Error() string {
	return e.msg
}

func (e *Err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}

func (e *Err) SetType(errorType ErrorType) *Err {
	e.errorType = errorType
	return e
}

func Wrap(err error, message string) *Err {
	return &Err{
		err: err,
		msg: fmt.Sprintf("%s: %s", message, err),
	}
}

func Wrapf(err error, message string, a ...interface{}) *Err {
	msg := fmt.Sprintf(message, a...)
	return &Err{
		err: err,
		msg: fmt.Sprintf("%s: %s", msg, err),
	}
}

func IsType(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*Err)
	if !ok {
		return errorType == ErrorTypeUnkown
	}

	if e.errorType == errorType {
		return true
	}

	return IsType(e.err, errorType)
}
