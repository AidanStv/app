package liberror

import "errors"

var (
	ErrUsersNotFound = errors.New("users not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserNotCreate = errors.New("user not create")
	ErrUserNotUpdate = errors.New("user not update")
	ErrUserNotDelete = errors.New("user not delete")
)

var (
	ErrProductssNotFound = errors.New("users not found")
	ErrProductsNotFound  = errors.New("user not found")
	ErrProductsNotCreate = errors.New("user not create")
	ErrProductsNotUpdate = errors.New("user not update")
	ErrProductsNotDelete = errors.New("user not delete")
)
