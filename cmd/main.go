package main

// @title subscriptions API
// @version 1.0
// @description API для агрегации записей о подписках

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/subscriptions_api/internal/config"
	"github.com/subscriptions_api/internal/logger"
	"github.com/subscriptions_api/internal/repository"
	"github.com/subscriptions_api/routes"
)

func main() {
	cfg := config.MustLoad()

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	logger.Init("text")

	/*	db := postgres.ConnectDB(cfg)
		postgres.CreateTables(db)
		defer db.Close(context.Background())*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Name)

	db, err := repository.NewPostgresDB(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to init DB", err)
	}
	defer db.Close(ctx)

	if err := repository.RunMigrations(ctx, db); err != nil {
		log.Fatal("migrations", err)
	}

	routes.InitRoutes(app)
	log.Fatal(app.Listen(cfg.Server.Port))
}
