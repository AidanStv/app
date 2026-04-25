package handler

import (
	"my-project/internal/model"
	"my-project/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserService *service.UserService
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
		return c.JSON(http.StatusBadRequest, err.Error())
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "user created")
}
