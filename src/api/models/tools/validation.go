package tools

import "gopkg.in/go-playground/validator.v8"

const (
	MaxTextLength          = 4096
	MaxRegularStringLength = 255
)

type ValidationErr struct {
	err error
}

func (e *ValidationErr) Error() string {
	if e == nil {
		return "null"
	}

	if e.err == nil {
		return "e.err = null"
	}
	return e.err.Error()
}

func NewValidationErr(err error) *ValidationErr {
	vErr := &ValidationErr{err: err}
	return vErr
}

type NotPermittedErr struct {
	err error
}

func (e *NotPermittedErr) Error() string {
	if e == nil {
		return "null"
	}

	if e.err == nil {
		return "e.err = null"
	}
	return e.err.Error()
}

func NewNotPermittedErr(err error) *NotPermittedErr {
	vErr := &NotPermittedErr{err: err}
	return vErr
}

var Validator = validator.New(&validator.Config{TagName: "validate"})
