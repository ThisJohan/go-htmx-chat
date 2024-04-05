package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ThisJohan/go-htmx-chat/handler"
	"github.com/ThisJohan/go-htmx-chat/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type config struct {
	port string
	psql models.PostgresConfig
}

func loadConfig() (cfg config, err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	cfg.port = os.Getenv("PORT")
	cfg.psql = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	return
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cfg config) error {
	e := echo.New()
	db, err := models.OpenDB(cfg.psql)
	if err != nil {
		return err
	}
	defer db.Close()

	userService := &models.UserService{
		DB: db,
	}

	userHandler := handler.UserHandler{
		UserService: userService,
	}

	e.Static("/assets", "assets")
	e.GET("/", userHandler.ShowUser)
	e.GET("/signup", userHandler.Signup)
	e.POST("/signup", userHandler.ProcessSignup)

	return e.Start(fmt.Sprintf(":%s", cfg.port))
}
