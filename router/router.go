package router

import (
	m "A-Simple-Api-Go-Fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"net/http"
	"time"
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
	req := c.GetReqHeaders()
	r := m.Pong{
		Agent:   req[fiber.HeaderUserAgent],
		Date:    time.Now().Format(http.TimeFormat),
		Message: "pong",
	}
	return c.JSON(r)
}

func info(c *fiber.Ctx) error {
	return c.SendString("info")
}

func connection(c *fiber.Ctx) error {
	return c.SendString("connection")
}
