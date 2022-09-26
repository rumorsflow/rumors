package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type UpdateHandler struct {
}

func (h *UpdateHandler) Register(r *echo.Group) {
	g := r.Group("/updates")
	{
		g.GET("", h.updates)
	}
}

func (h *UpdateHandler) updates(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
