package main

import (
	"context"
	"fmt"
	"my-project/internal/handler"
	"my-project/internal/repository"
	"my-project/internal/service"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

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

	e.GET("/users", h.GetUsers)
	// e.GET("/users/:id", h.GetUser)
	e.POST("/users", h.CreateUser)
	e.PATCH("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteHandler)

	e.Start(":8080")
}
