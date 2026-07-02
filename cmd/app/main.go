package main

import (
	"context"
	"fmt"
	"my-project/internal/handler"
	"my-project/internal/middleware"
	"my-project/internal/repository"
	"my-project/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// TODO:
// 1) Logout
// 2) /users --->>> /products (новая структура таблицы)
func main() {

	ctx := context.Background()
	url := "postgres://postgres:postgres@localhost:5432/popcorn"
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Println("Unable to connect to database", err)
		return
	}
	defer conn.Close(ctx)

	UserRepository := &repository.UserRepository{
		Conn: conn,
	}

	ProductRepository := &repository.ProductRepository{
		Conn: conn,
	}

	UserService := service.NewUserService(UserRepository)

	ProductService := &service.ProductService{
		ProductRepository: ProductRepository,
	}

	h := &handler.Handler{
		UserService:    UserService,
		ProductService: ProductService,
	}

	e := echo.New()

	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
	e.POST("/refresh", h.Refresh)

	// CRUD Products
	e.POST("/products", h.CreateProduct, middleware.JWTMiddleware)
	e.GET("/products", h.GetProducts, middleware.JWTMiddleware)
	e.GET("/products/:id", h.GetProduct, middleware.JWTMiddleware)
	e.PATCH("/products/:id", h.UpdateProduct, middleware.JWTMiddleware)
	e.DELETE("/products/:id", h.DeleteProduct, middleware.JWTMiddleware)

	// CRUD Users
	e.POST("/users", h.CreateUser, middleware.JWTMiddleware)
	e.GET("/users", h.GetUsers, middleware.JWTMiddleware)
	e.GET("/users/:id", h.GetUser, middleware.JWTMiddleware)
	e.PATCH("/users/:id", h.UpdateUser, middleware.JWTMiddleware)
	e.DELETE("/users/:id", h.DeleteHandler, middleware.JWTMiddleware)
	e.POST("/logout", h.Logout, middleware.JWTMiddleware)

	e.Start(":8080")
}
