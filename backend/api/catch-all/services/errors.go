package services

import "errors"

var (
	ErrUnexpected          = errors.New("unexpected server error")
	ErrThumbnailNotCreated = errors.New("thumbnail has not been created")
)
