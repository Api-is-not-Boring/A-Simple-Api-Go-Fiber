package router

import (
	m "A-Simple-Api-Go-Fiber/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var db, _ = m.DbInit()

// CreateRoutes for fiber app
func CreateRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// Api v1
	v1 := api.Group("/v1")
	v1.Get("/ping", ping)
	v1.Get("/info", info)
	v1.Get("/connections", connections)

	// Api v2 Restful
	v2 := api.Group("/v2")
	v2.Get("/cars", getCars)
	v2.Post("/cars", createCar)
	v2.Put("/cars", updateCar)
	v2.Delete("/cars", deleteCar)
	v2.Get("/cars/:id<int>", getCars)
	v2.Put("/cars/:id<int>", updateCar)
	v2.Delete("/cars/:id<int>", deleteCar)

	// TODO Api v3 Auth Jwt
	v3 := api.Group("/v3")
	v3.Get("/auth", announce)

}

// Endpoint Api v1
func ping(c *fiber.Ctx) error {
	return c.JSON(m.GetPong(c))
}

func info(c *fiber.Ctx) error {
	return c.JSON(m.GetInfo())
}

func connections(c *fiber.Ctx) error {
	return c.JSON(m.GetConnections())
}

// GetCars Endpoint Api v2
func getCars(c *fiber.Ctx) error {
	return c.JSON(m.GetCars(c, db))
}

func createCar(c *fiber.Ctx) error {
	return c.JSON(m.CreateCar(c, db))
}

func updateCar(c *fiber.Ctx) error {
	return c.JSON(m.UpdateCar(c, db))
}

func deleteCar(c *fiber.Ctx) error {
	return c.JSON(m.DeleteCar(c, db))
}

func announce(c *fiber.Ctx) error {
	return c.JSON(m.Announce(c))
}
