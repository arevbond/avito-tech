package utils

import (
	"database/sql"
	"errors"
	"net/http"
)

const (
	notFoundMessage      string = "not found"
	notExistUserMessage  string = "user doesn't exist"
	InternalErrorMessage string = "internal error"
)

type ErrorResult struct {
	Err        error
	Msg        string
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return e.Err.Error()
}

func WrapInternalError(err error) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        InternalErrorMessage,
		StatusCode: http.StatusInternalServerError,
	}
}

func WrapError(err error, msg string, status int) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: status,
	}
}

func WrapSqlError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return WrapError(err, notFoundMessage, http.StatusNotFound)
	default:
		return WrapInternalError(err)
	}
}

func FromError(err error) (*ErrorResult, bool) {
	if err == nil {
		return nil, false
	}

	var result *ErrorResult
	ok := errors.As(err, &result)
	if !ok {
		return nil, false
	}
	return result, true
}
