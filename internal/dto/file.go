package dto

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrFileDeleteFailed = errors.New("failed to delete file")
)

const (
	MESSAGE_FILE_FETCH_SUCCESS = "File fetched successfully"
	MESSAGE_FILE_FETCH_FAILED  = "Failed to fetch file"
)
