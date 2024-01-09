package dto

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrFileDeleteFailed = errors.New("failed to delete file")
)

const (
	MessageFileFetchSuccess = "File fetched successfully"
	MessageFileFetchFailed  = "Failed to fetch file"
)
