package router

import (
	m "A-Simple-Api-Go-Fiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm/clause"
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
	v2.Get("/cars/:id<int>", getCarWithPath)
	v2.Put("/cars/:id<int>", updateCarWithPath)
	v2.Delete("/cars/:id<int>", deleteCarWithPath)
	// TODO Api v3 Auth Jwt
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
	id := c.Query("id")
	if id == "" {
		var cars []m.Car
		db.Find(&cars)
		return c.JSON(m.AllCars{Total: len(cars), Cars: cars})
	} else {
		var car m.Car
		db.First(&car, id)
		return c.JSON(car)
	}
}

func createCar(c *fiber.Ctx) error {
	var car m.Car
	var records int64
	if err := c.BodyParser(&car); err != nil {
		return err
	}
	db.Model(&[]m.Car{}).Count(&records)
	if records >= 20 {
		m.DbReset(db)
	}
	db.Create(&car)
	return c.JSON(car)
}

func updateCar(c *fiber.Ctx) error {
	id := c.Query("id")
	var car m.Car
	if err := c.BodyParser(&car); err != nil {
		return err
	}
	if id == "" {
		db.Model(&car).Updates(car)
		return c.JSON(car)
	} else {
		db.Model(&car).Clauses(clause.Returning{}).Where("id = ?", id).Updates(car)
		return c.JSON(car)
	}
}

func deleteCar(c *fiber.Ctx) error {
	id := c.Query("id")
	var car m.Car
	db.First(&car, id)
	db.Delete(&car)
	return c.JSON(car)
}

func getCarWithPath(c *fiber.Ctx) error {
	id := c.Params("id")
	var car m.Car
	db.First(&car, id)
	return c.JSON(car)
}

func updateCarWithPath(c *fiber.Ctx) error {
	id := c.Params("id")
	var car m.Car
	if err := c.BodyParser(&car); err != nil {
		return err
	}
	db.Model(&car).Clauses(clause.Returning{}).Where("id = ?", id).Updates(car)
	return c.JSON(car)
}

func deleteCarWithPath(c *fiber.Ctx) error {
	id := c.Params("id")
	var car m.Car
	db.First(&car, id)
	db.Delete(&car)
	return c.JSON(car)
}
