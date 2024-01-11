package errors

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserNoPicture      = errors.New("user don't have any picture")
)
