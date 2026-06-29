package model

import (
	"errors"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	University   string `json:"university"`
}

type RegisterRequest struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	University string `json:"university"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID     int
	UserID int
	Token  string
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
