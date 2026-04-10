package api

import "github.com/labstack/echo/v5"

type MainHandler struct {
	GameApiHandler *GameApiHandler
}

func NewMainHandler() *MainHandler {
	return &MainHandler{
		GameApiHandler: &GameApiHandler{},
	}
}

func (h *MainHandler) GetGame(ctx *echo.Context) error {
	return h.GameApiHandler.GetGame(ctx)
}
