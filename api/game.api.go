package api

import "github.com/labstack/echo/v5"

type GameApiHandler struct {
}

func (h *GameApiHandler) GetGame(ctx *echo.Context) error {
	return ctx.JSON(200, []Character{})
}
