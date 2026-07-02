package handler

import (
	"errors"
	"my-project/internal/model"
	"my-project/pkg/liberror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetProduct(c echo.Context) error {

	ctx := c.Request().Context()

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid product id")
	}

	product, err := h.ProductService.GetProduct(ctx, id)
	if err != nil {

		if errors.Is(err, liberror.ErrProductsNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		}

		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProducts(c echo.Context) error {

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

	products, err := h.ProductService.GetProducts(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, liberror.ErrProductsNotFound) {
			return c.JSON(http.StatusNotFound, "Products not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}

	ctx := c.Request().Context()

	err = h.ProductService.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, liberror.ErrProductsNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "Product deleted")
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad id")
	}

	var p model.Product
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	p.ID = id

	err = h.ProductService.UpdateProduct(c.Request().Context(), p)
	if err != nil {
		if errors.Is(err, liberror.ErrProductsNotFound) {
			return c.JSON(http.StatusNotFound, "Product not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, "Product updated")
}

func (h *Handler) CreateProduct(c echo.Context) error {

	var p model.Product

	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	err := h.ProductService.CreateProduct(c.Request().Context(), p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Product created")
}
