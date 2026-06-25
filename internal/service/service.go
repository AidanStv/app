package service

import (
	"context"
	"errors"
	"my-project/internal/model"
	"my-project/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) Register(ctx context.Context, req model.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	user := model.User{
		Name:         req.Name,
		University:   req.University,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}
	return s.UserRepository.CreateUser(ctx, user)
}

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

func (s *UserService) GetUsers(ctx context.Context, limit, offset int) ([]model.User, error) {
	return s.UserRepository.GetUsers(ctx, limit, offset)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.UserRepository.GetByEmail(ctx, email)
}

func (s *UserService) GetUser(ctx context.Context, id int) (model.User, error) {
	if id <= 0 {
		return model.User{}, errors.New("invalid id")
	}

	return s.UserRepository.GetUser(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, u model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	return s.UserRepository.CreateUser(ctx, u)
}

func (s *UserService) UpdateUser(ctx context.Context, u model.User) error {
	if u.ID <= 0 {
		return errors.New("invalid id")
	}

	if err := u.Validate(); err != nil {
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
