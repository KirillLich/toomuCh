package tmerrors

import (
	"errors"
	"net/http"
)

var ErrTooLongText = errors.New("too long text")
var ErrTooLongTitle = errors.New("too long title")
var ErrInvalidLimit = errors.New("invalid limit value")
var ErrInvalidId = errors.New("invalid id value")
var ErrInvalidTime = errors.New("invalid time value")

// var ErrEmptyTitle = errors.New("empty title of todo")
// var ErrNotFound = errors.New("todo not found")

func MapErrorHttpStatus(err error) (int, string) {
	if errors.Is(err, ErrTooLongText) {
		return http.StatusBadRequest, "too long text"
	}
	if errors.Is(err, ErrTooLongTitle) {
		return http.StatusBadRequest, "too long title"
	}
	if errors.Is(err, ErrInvalidId) {
		return http.StatusBadRequest, "invalid id"
	}
	if errors.Is(err, ErrInvalidLimit) {
		return http.StatusBadRequest, "invalid limit"
	}
	if errors.Is(err, ErrInvalidTime) {
		return http.StatusBadRequest, "invalid time"
	}
	return http.StatusInternalServerError, "internal server error"
}
