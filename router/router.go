package router

import (
	m "A-Simple-Api-Go-Fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// CreateRoutes for fiber app
func CreateRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// Api v1
	v1 := api.Group("/v1")
	v1.Get("/ping", ping)
	v1.Get("/info", info)
	v1.Get("/connection", connection)

	// TODO Api v2 Restful

	// TODO Api v3 Auth Jwt
}

// Ping endpoint Api v1
func ping(c *fiber.Ctx) error {
	return c.JSON(m.GetPong(c))
}

func info(c *fiber.Ctx) error {
	return c.JSON(m.GetInfo())
}

func connection(c *fiber.Ctx) error {
	return c.SendString("connection")
}
