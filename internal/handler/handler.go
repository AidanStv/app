package handler

import (
	"errors"
	"my-project/internal/model"
	"my-project/internal/service"
	"my-project/pkg/liberror"

	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Register(c echo.Context) error {
	var req model.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	err := h.UserService.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, "user created")
}

func (h *Handler) Login(c echo.Context) error {
	var req model.LoginRequest

	//с.Bind записывает из боди в req
	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := h.UserService.GetEmail(
		c.Request().Context(),
		req.Email,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"user not found",
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"invalid password",
		)
	}
	return c.JSON(200, "login success")
}

type Handler struct {
	UserService *service.UserService
}

func (h *Handler) GetUser(c echo.Context) error {

	ctx := c.Request().Context()

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid user id")
	}

	user, err := h.UserService.GetUser(ctx, id)
	if err != nil {

		if errors.Is(err, liberror.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, "user not found")
		}

		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUsers(c echo.Context) error {

	ctx := c.Request().Context()

	limit := 10
	page := 1

	if p := c.QueryParam("page"); p != "" {
		val, err := strconv.Atoi(p)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "invalid page")
		}
		page = val
	}

	offset := (page - 1) * limit

	users, err := h.UserService.GetUsers(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, liberror.ErrUsersNotFound) {
			return c.JSON(http.StatusNotFound, "users not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteHandler(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}

	ctx := c.Request().Context()

	err = h.UserService.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, liberror.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, "user not found")
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "user deleted")
}

func (h *Handler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}

	var u model.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	u.ID = id

	err = h.UserService.UpdateUser(c.Request().Context(), u)
	if err != nil {
		if errors.Is(err, liberror.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, "user updated")
}

func (h *Handler) CreateUser(c echo.Context) error {

	var u model.User

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	err := h.UserService.CreateUser(c.Request().Context(), u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "user created")
}
