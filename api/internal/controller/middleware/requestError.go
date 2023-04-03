package middleware

import "errors"

type requestError struct {
	StatusCode int
	Err        error
}

func (r *requestError) Error() string {
	return r.Err.Error()
}

func NewStr(statusCode int, err string) error {
	return &requestError{
		StatusCode: statusCode,
		Err:        errors.New(err),
	}
}

func NewErr(statusCode int, err error) error {
	return &requestError{
		StatusCode: statusCode,
		Err:        err,
	}
}
