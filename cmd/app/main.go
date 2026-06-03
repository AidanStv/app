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
// 1) errors.Is()
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

	UserService := service.NewUserService(UserRepository)

	h := &handler.Handler{UserService: UserService}

	e := echo.New()

	e.POST("/login", h.Login)
	e.POST("/register", h.Register)

	e.GET("/users", h.GetUsers, middleware.JWTMiddleware)
	e.GET("/users/:id", h.GetUser, middleware.JWTMiddleware)
	e.POST("/users", h.CreateUser, middleware.JWTMiddleware)
	e.PATCH("/users/:id", h.UpdateUser, middleware.JWTMiddleware)
	e.DELETE("/users/:id", h.DeleteHandler, middleware.JWTMiddleware)

	e.Start(":8080")
}
