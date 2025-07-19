package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/subscriptions_api/docs"
	"github.com/subscriptions_api/handlers"
)

func InitRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/subscriptions", handlers.CreateSubscription)
	api.Get("/subscriptions/:id", handlers.GetSubscription)
	api.Put("/subscriptions/:id", handlers.UpdateSubscription)
	api.Delete("/subscriptions/:id", handlers.DeleteSubscription)
	api.Get("/subscriptions", handlers.GetAllSubscriptions)
	api.Get("/total", handlers.GetTotalPriceInPeriod)
	app.Get("/swagger/*", swagger.HandlerDefault)
}
