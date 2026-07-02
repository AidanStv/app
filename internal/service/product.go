package service

import (
	"context"
	"errors"
	"my-project/internal/model"
	"my-project/internal/repository"
)

type ProductService struct {
	ProductRepository *repository.ProductRepository
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]model.Product, error) {
	return s.ProductRepository.GetProducts(ctx, limit, offset)
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (model.Product, error) {
	if id <= 0 {
		return model.Product{}, errors.New("invalid id")
	}

	return s.ProductRepository.GetProduct(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, p model.Product) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return s.ProductRepository.CreateProduct(ctx, p)
}

func (s *ProductService) UpdateProduct(ctx context.Context, p model.Product) error {
	if p.ID <= 0 {
		return errors.New("invalid id")
	}

	if err := p.Validate(); err != nil {
		return err
	}

	return s.ProductRepository.UpdateProduct(ctx, p)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	return s.ProductRepository.DeleteProduct(ctx, id)
}
