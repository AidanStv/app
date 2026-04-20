package service

import (
	"context"
	"errors"
	"my-project/internal/model"
	"my-project/internal/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{UserRepository: r}
}

type Pagination struct {
	Page  int
	Limit int
}

func (s *UserService) validateID(u model.User) error {
	if u.ID <= 0 {
		return errors.New("invalid user id")
	}

	if u.Name == "" {
		return errors.New("invalid name")
	}

	if u.University == "" {
		return errors.New("invalid university")
	}
	return nil
}

func (s *UserService) validate(u model.User) error {
	if u.Name == "" {
		return errors.New("invalid name")
	}

	if u.University == "" {
		return errors.New("invalid university")
	}

	return nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	return s.UserRepository.GetUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, errors.New("invalid id")
	}

	return s.UserRepository.GetUser(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, u model.User) error {
	if err := s.validateID(u); err != nil {
		return err
	}

	return s.UserRepository.CreateUser(ctx, u)
}

func (s *UserService) UpdateUser(ctx context.Context, u model.User) error {
	if err := s.validate(u); err != nil {
		return err
	}

	return s.UserRepository.UpdateUser(ctx, u)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	return s.UserRepository.DeleteUser(ctx, id)
}
