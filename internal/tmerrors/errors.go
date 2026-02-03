package tmerrors

import (
	"errors"
)

var ErrTooLongText = errors.New("too long text")
var ErrTooLongTitle = errors.New("too long title")
var ErrInvalidLimit = errors.New("invalid limit value")
var ErrInvalidId = errors.New("invalid id value")

// var ErrEmptyTitle = errors.New("empty title of todo")
// var ErrNotFound = errors.New("todo not found")

// func MapErrorHttpStatus(err error) (int, string) {
// 	if errors.Is(err, ErrEmptyTitle) {
// 		return http.StatusBadRequest, "empty title"
// 	}
// 	if errors.Is(err, ErrInvalidId) {
// 		return http.StatusBadRequest, "invalid id"
// 	}
// 	if errors.Is(err, ErrNotFound) {
// 		return http.StatusNotFound, "not found"
// 	}
// 	return http.StatusInternalServerError, "internal server error"
// }
