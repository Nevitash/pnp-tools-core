package main

import (
	"fmt"
	"gg/nevitash/pnp/tools/core/api"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	fmt.Println("This is a start!")
	echo, err := initEcho()
	if err != nil {
		log.Fatal("Failed to initialize Echo: ", err)
	}
	if err := echo.Start(":8080"); err != nil {
		log.Fatal("Server crashed: ", err)
	}
}

func initEcho() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.RequestLogger(), middleware.RequestID())
	mainHandler := api.NewMainHandler()
	api.RegisterHandlers(e, mainHandler)
	log.Println("Starting server on http://localhost:8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Server crashed: ", err)
		return nil, err
	}
	return e, nil
}
