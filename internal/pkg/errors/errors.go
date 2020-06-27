package errors

import (
	"fmt"
	"io"
)

// ErrorType defines an errors behavior
type ErrorType int

const (
	// ErrorTypeUnknown default error
	ErrorTypeUnknown ErrorType = iota

	// ErrorTypeBadRequest for bad input error
	ErrorTypeBadRequest

	// ErrorTypeNotFound for resources that are not existent
	ErrorTypeNotFound

	// ErrorTypeInternal for internal error
	ErrorTypeInternal
)

// Err represents a single error
type Err struct {
	err       error
	errorType ErrorType
	msg       string
}

// New creates a new error
func New(s string) *Err {
	return &Err{
		msg: s,
	}
}

// Errorf creates a new error using a formatted string
func Errorf(s string, a ...interface{}) *Err {
	return &Err{
		msg: fmt.Sprintf(s, a...),
	}
}

// Error returns the message from within an error
func (e *Err) Error() string {
	return e.msg
}

// Format the errors according to the verbs
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

// SetType set the error type
func (e *Err) SetType(errorType ErrorType) *Err {
	e.errorType = errorType
	return e
}

// Wrap an existing error in more contextual information
func Wrap(err error, message string) *Err {
	return &Err{
		err: err,
		msg: fmt.Sprintf("%s: %s", message, err),
	}
}

// Wrapf an existing error in more contextual information and allow for format operations
func Wrapf(err error, message string, a ...interface{}) *Err {
	msg := fmt.Sprintf(message, a...)
	return &Err{
		err: err,
		msg: fmt.Sprintf("%s: %s", msg, err),
	}
}

// IsType checks whether an error is of a given Type
func IsType(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*Err)
	if !ok {
		return errorType == ErrorTypeUnknown
	}

	if e.errorType == errorType {
		return true
	}

	return IsType(e.err, errorType)
}
