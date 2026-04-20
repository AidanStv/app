package repository

import (
	"context"
	"errors"
	"fmt"
	"my-project/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Conn *pgx.Conn
}

// TODO:
// 1) Rename error messages (only english) +++
// 2) paste nil, not empty slice +++
// 3) in main.go check errors. +++
// 4) in rows.Affected use errors.New() +++
// 5) text of errors starts with lowercase letter +++
// 6) in r.DeleteUser use only id, not full user struct +++

func (r *UserRepository) GetUsers(ctx context.Context) ([]model.User, error) {

	rows, err := r.Conn.Query(ctx, "select id, name, university from users")
	if err != nil {
		return nil, fmt.Errorf("SQL query execution: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Name, &u.University)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return users, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (model.User, error) {
	var u model.User
	err := r.Conn.QueryRow(ctx, "select id, name, university from users where id = $1", id).Scan(&u.ID, &u.Name, &u.University)
	if err != nil {
		return model.User{}, fmt.Errorf("scan user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u model.User) error {
	commandTag, err := r.Conn.Exec(context.Background(), "insert into users(name, university) values($1, $2)", u.Name, u.University)
	if err != nil {
		return fmt.Errorf("SQL query execution error: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("user not created")
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
		return errors.New("user not updated")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	commandTag, err := r.Conn.Exec(context.Background(), "delete from users where id = $1", id)
	if err != nil {
		return fmt.Errorf("SQL query execution: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("user not deleted")
	}
	return nil
}
