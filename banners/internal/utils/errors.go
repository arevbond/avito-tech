package utils

import (
	"banners/internal/service"
	"database/sql"
	"errors"
	"net/http"
)

const (
	NotFoundMessage      string = "not found"
	InternalErrorMessage string = "internal error"
)

var (
	ErrInvalidTypeParam = &ErrorResult{Err: errors.New("invalid type param"),
		Msg: "invalid type param", StatusCode: http.StatusBadRequest}
	ErrNotRequiredParam = &ErrorResult{Err: errors.New("not required param"),
		Msg: "required param is absent", StatusCode: http.StatusBadRequest}
)

type ErrorResult struct {
	Err        error
	Msg        string
	StatusCode int
}

func (e *ErrorResult) Error() string {
	return e.Err.Error()
}

func WrapForbiddenError(err error, msg string) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: http.StatusForbidden,
	}
}

func WrapNotFoundError(err error, msg string) *ErrorResult {
	return &ErrorResult{
		Err:        err,
		Msg:        msg,
		StatusCode: http.StatusNotFound,
	}
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

func WrapServiceError(err error) *ErrorResult {
	switch err {
	case service.ErrUnauthorized:
		return &ErrorResult{
			Err:        err,
			Msg:        err.Error(),
			StatusCode: 401,
		}
	case service.ErrNotFound:
		return &ErrorResult{
			Err:        err,
			Msg:        err.Error(),
			StatusCode: 404,
		}
	case service.ErrForbidden:
		return &ErrorResult{
			Err:        err,
			Msg:        err.Error(),
			StatusCode: 403,
		}
	case service.ErrIncorrectData:
		return &ErrorResult{
			Err:        err,
			Msg:        err.Error(),
			StatusCode: 400,
		}
	default:
		return WrapInternalError(err)
	}
}

func WrapSqlError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return WrapNotFoundError(err, NotFoundMessage)
	default:
		return WrapInternalError(err)
	}
}
