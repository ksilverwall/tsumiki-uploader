package services

import "errors"

type ErrorCode int

const (
	ErrLabelUnexpected ErrorCode = iota
	ErrLabelThumbnailNotCreated
)

type Error struct {
	Code   ErrorCode
	message string
	cause   error
}

func (e *Error) Error() string {
	return e.message
}

func NewError(label ErrorCode, cause error) Error {
	return Error{
		Code:   label,
		message: cause.Error(),
	}
}

var (
	ErrUnexpected          = errors.New("unexpected server error")
	ErrThumbnailNotCreated = errors.New("thumbnail has not been created")
)
