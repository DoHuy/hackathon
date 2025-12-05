package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"hackathon/config"
	"hackathon/database"
	"hackathon/handlers"
	"hackathon/middleware"
	"hackathon/pkg/logger"
	customValidator "hackathon/pkg/validator"
	"hackathon/repositories"
	"hackathon/services"

	_ "hackathon/docs"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Hackathon API
// @version 1.0
// @description API Server with JWT Auth, Postgres
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.Server.LogLevel, cfg.Server.PrettyLog)
	log.Info().Msg("Starting application...")

	database.Init(cfg.Database)
	repos := repositories.NewRepository(database.DB)

	allowedTypes := strings.Split(cfg.Storage.AllowedTypes, ",")
	srv := services.NewService(repos, []byte(cfg.JWT.Secret), cfg.JWT.ExpirationHours, cfg.Storage.UploadDir, cfg.Storage.MaxSizeMB, allowedTypes)

	e := echo.New()
	e.Validator = customValidator.NewCustomValidator()
	e.Use(middleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	handlers.NewHandler(e.Group("/api"), srv, cfg, repos).RegisterRoutes()

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Server started")
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server crashed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Shutdown error")
	}

	database.Close()
	log.Info().Msg("Server exited properly")
}
