package main

import (
	"context"
	"database/sql"
	"fmt"
	"gg/nevitash/pnp/tools/core/ent"
	"gg/nevitash/pnp/tools/core/api"
	"log"
	"net/url"
	"os"
	"strings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "gg/nevitash/pnp/tools/core/resources/db"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	fmt.Println("This is a start!")
	client, err := initEnt()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer client.Close()

	e, err := initEcho()
	if err != nil {
		log.Fatal("Failed to initialize Echo: ", err)
	}
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Server crashed: ", err)
	}
}

func initEcho() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.RequestLogger(), middleware.RequestID())
	mainHandler := api.NewMainHandler()
	api.RegisterHandlers(e, mainHandler)
	log.Println("Starting server on http://localhost:8080")
	return e, nil
}

func initEnt() (*ent.Client, error) {
	dsn := determineDSN()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func determineDSN() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		return dsn
	}

	gooseDSN := strings.Trim(os.Getenv("GOOSE_DBSTRING"), "\"")
	if gooseDSN != "" {
		return gooseDSN
	}

	user := getenvDefault("POSTGRES_USER", "postgres")
	password := getenvDefault("POSTGRES_PASSWORD", "postgres")
	host := getenvDefault("POSTGRES_HOSTNAME", "localhost")
	dbName := getenvDefault("POSTGRES_DB", "postgres")
	port := getenvDefault("POSTGRES_PORT", "5432")

	dsnURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + port,
		Path:   "/" + dbName,
	}
	query := dsnURL.Query()
	query.Set("sslmode", getenvDefault("POSTGRES_SSLMODE", "disable"))
	dsnURL.RawQuery = query.Encode()
	return dsnURL.String()
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
