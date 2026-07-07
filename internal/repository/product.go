package repository

import (
	"context"
	"errors"
	"fmt"
	"my-project/internal/model"
	"my-project/pkg/liberror"

	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	Conn *pgx.Conn
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit, offset int) ([]model.Product, error) {

	rows, err := r.Conn.Query(ctx, "SELECT id, name, price FROM products ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("SQL query execution: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, fmt.Errorf("scan products: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	if len(products) == 0 {
		return nil, liberror.ErrProductsNotFound
	}

	return products, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int) (model.Product, error) {
	var p model.Product

	err := r.Conn.QueryRow(ctx, "select id, name, price from products where id = $1", id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, liberror.ErrProductsNotFound
		}
		return model.Product{}, fmt.Errorf("scan product: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p model.Product) error {
	commandTag, err := r.Conn.Exec(ctx, "insert into products(name, price) values($1, $2)", p.Name, p.Price)
	if err != nil {
		return fmt.Errorf("SQL query execution error: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrProductsNotCreate
	}
	return nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, p model.Product) error {
	query := "UPDATE products SET name = $1, price = $2 WHERE id = $3"
	commandTag, err := r.Conn.Exec(context.Background(), query, p.Name, p.Price, p.ID)
	if err != nil {
		return fmt.Errorf("SQL query execution error: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrProductsNotUpdate
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	commandTag, err := r.Conn.Exec(context.Background(), "delete from products where id = $1", id)
	if err != nil {
		return fmt.Errorf("SQL query execution: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return liberror.ErrProductsNotDelete
	}
	return nil
}
