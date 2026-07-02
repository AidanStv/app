package model

import "errors"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProduct struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	return nil
}
