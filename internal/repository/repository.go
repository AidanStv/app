package repository

import (
	"context"
	"errors"
	"fmt"
	"my-project/internal/model"
	"my-project/pkg/liberror"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Conn *pgx.Conn
}

func (r *UserRepository) GetUsers(ctx context.Context, limit, offset int) ([]model.User, error) {

	rows, err := r.Conn.Query(ctx, "SELECT id, name, university, email FROM users ORDER BY id limit $1 offset $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SQL query execution: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.University, &u.Email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	if len(users) == 0 {
		return nil, liberror.ErrUsersNotFound
	}

	return users, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.Conn.QueryRow(ctx, "SELECT id, email, password_hash FROM users WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {

		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (model.User, error) {
	var u model.User
	err := r.Conn.QueryRow(ctx, "select id, email, name, university from users where id = $1", id).Scan(&u.ID, &u.Email, &u.Name, &u.University)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, liberror.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("scan user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u model.User) error {
	commandTag, err := r.Conn.Exec(context.Background(), "insert into users(name, university, email, password_hash) values($1, $2, $3, $4)", u.Name, u.University, u.Email, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("SQL query execution error: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrUserNotCreate
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u model.User) error {
	query := "update users set name = $1, university = $2 where id = $3"
	commandTag, err := r.Conn.Exec(context.Background(), query, u.Name, u.University, u.ID)
	if err != nil {
		return fmt.Errorf("SQL query execution error: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrUserNotUpdate
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	commandTag, err := r.Conn.Exec(context.Background(), "delete from users where id = $1", id)
	if err != nil {
		return fmt.Errorf("SQL query execution: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrUserNotDelete
	}
	return nil
}
