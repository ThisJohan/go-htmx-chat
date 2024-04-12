package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ThisJohan/go-htmx-chat/handler"
	"github.com/ThisJohan/go-htmx-chat/models"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
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
	err = run(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config) error {
	e := echo.New()
	db, err := models.OpenDB(cfg.psql)
	if err != nil {
		return err
	}
	defer db.Close()

	redis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	err = redis.Ping(ctx).Err()
	if err != nil {
		return err
	}

	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		Redis: redis,
	}
	chatService := models.NewChatService(db, redis)
	go chatService.Hub.Run()

	userHandler := handler.UserHandler{
		UserService:    userService,
		SessionService: sessionService,
	}
	socketHandler := handler.SocketHandler{
		ChatService: chatService,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", func(c echo.Context) error {
		cl := chatService.Hub.Register(nil, 1000)
		m := &models.Message{
			ToUser:  1,
			Content: "Tesst",
		}
		cl.ReadPump(m)
		return c.String(200, "ok")
	})

	e.Static("/assets", "assets")
	e.GET("/", userHandler.ShowUser)
	e.GET("/signup", userHandler.Signup)
	e.POST("/signup", userHandler.ProcessSignup)
	e.GET("/login", userHandler.Login)
	e.POST("/login", userHandler.ProcessLogin)

	g := e.Group("/user")
	g.Use(userHandler.AuthRequired)
	g.GET("/me", userHandler.Me)

	e.GET("/chat", socketHandler.Demo)
	g.GET("/ws", socketHandler.Chat)

	return e.Start(fmt.Sprintf(":%s", cfg.port))
}
