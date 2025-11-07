package errs

import (
	"errors"
	"fmt"
)

type Code string

const (
	InvalidInput Code = "INVALID_INPUT"
	NotFound     Code = "NOT_FOUND"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	Conflict     Code = "CONFLICT"
	Internal     Code = "INTERNAL"
	Unavailable  Code = "UNAVAILABLE"
	Timeout      Code = "TIMEOUT"
)

type Op string

type E struct {
	Op     Op
	Code   Code
	Err    error
	Fields map[string]any
}

func (e *E) Error() string {
	switch {
	case e.Op != "" && e.Code != "":
		return fmt.Sprintf("%s %s: %v", e.Op, e.Code, e.Err)
	case e.Code != "":
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	default:
		return fmt.Sprintf("%v", e.Err)
	}
}

func (e *E) Unwrap() error { return e.Err }

func New(op Op, code Code, err error, fields map[string]any) *E {
	if err == nil {
		return nil
	}
	return &E{Op: op, Code: code, Err: err, Fields: fields}
}

func Wrap(op Op, err error, code Code) error {
	if err == nil {
		return nil
	}
	var e *E
	if errors.As(err, &e) {
		if e.Op == "" {
			e.Op = op
		}
		return e
	}
	return &E{Op: op, Code: code, Err: err}
}

func IsCode(err error, code Code) bool {
	var e *E
	return errors.As(err, &e) && e.Code == code
}
