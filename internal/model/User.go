package model

import (
	"errors"
)

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	University string `json:"university"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("invalid name")
	}

	if u.University == "" {
		return errors.New("invalid university")
	}

	return nil
}
