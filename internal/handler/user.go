package handler

import (
	"errors"
	"my-project/internal/model"
	"my-project/internal/service"
	"my-project/pkg/jwt"
	"my-project/pkg/liberror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserService    *service.UserService
	ProductService *service.ProductService
}

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

	user, err := h.UserService.GetByEmail(c.Request().Context(), req.Email)
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
			"bad password",
		)
	}

	accessToken, err := jwt.GenerateAccessToken(
		user.ID,
		user.Email,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	refreshToken, err := jwt.GenerateRefreshToken(
		user.ID,
		user.Email,
	)

	err = h.UserService.SaveRefreshToken(c.Request().Context(), user.ID, refreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	var req model.RefreshRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	claims, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return err
	}

	exists, err := h.UserService.RefreshTokenExists(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusUnauthorized, "refresh token not found")
	}

	accessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
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

func (h *Handler) Logout(c echo.Context) error {

	var req model.RefreshRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	err := h.UserService.DeleteRefreshToken(
		c.Request().Context(),
		req.RefreshToken,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "logout success")
}
