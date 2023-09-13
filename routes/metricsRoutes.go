package routes

import (
	"cpu-mon/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(server *fiber.App) {
	metrics := server.Group("/metrics")
	metrics.Get("/cpu/ssh", handlers.RunSsh)
	metrics.Get("/users", handlers.FetchAllUsers)
	metrics.Get("/users/username", handlers.FetchUser)
}
