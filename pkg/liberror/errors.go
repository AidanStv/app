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
	ErrProductssNotFound = errors.New("product not found")
	ErrProductsNotFound  = errors.New("product not found")
	ErrProductsNotCreate = errors.New("product not create")
	ErrProductsNotUpdate = errors.New("product not update")
	ErrProductsNotDelete = errors.New("product not delete")
)
