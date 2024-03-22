package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Armunz/learn-mongodb/internal/config"
	"github.com/Armunz/learn-mongodb/internal/controllers"
	"github.com/Armunz/learn-mongodb/internal/repositories"
	"github.com/Armunz/learn-mongodb/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	validate := validator.New()
	cfg := config.New(validate)

	mongoDB := config.NewMongo(ctx, cfg)

	// init repo
	repo := repositories.New(mongoDB, cfg.AppMongoQueryTimeoutMs)

	// init service
	service := services.NewService(repo, cfg.DefaultLimit)

	// init fiber
	app := fiber.New()
	app.Use(
		recover.New(),
		cors.New(cors.Config{
			AllowHeaders: "*",
		}),
	)

	// init controller
	controllers.RegisterHandlers(app.Group("/accounts"), service, validate, cfg.APITimeout)

	// Listen from a different goroutine
	address := ":9999"
	go func() {
		if err := app.Listen(address); err != nil {
			log.Err(err).Caller().Msg("failed to serve fiber http server")
			os.Exit(-1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	// close database
	log.Info().Msg("Closing MongoDB Connection...")
	if err := mongoDB.Client().Disconnect(ctx); err != nil {
		log.Err(err).Caller().Msg("failed to close MySQL database")
	}

	// close fiber
	log.Info().Msg("Shuting down Fiber server...")
	if err := app.Shutdown(); err != nil {
		log.Err(err).Caller().Msg("failed to shutdown fiber server")
	}
}
