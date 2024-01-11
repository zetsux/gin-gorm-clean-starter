package errors

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrFileDeleteFailed = errors.New("failed to delete file")
)
