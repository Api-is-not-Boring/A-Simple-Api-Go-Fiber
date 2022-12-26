package main

import (
	"A-Simple-Api-Go-Fiber/router"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"strconv"
)

func main() {
	// Load .env file
	_ = godotenv.Load(".env")
	listeningAddress := os.Getenv("ADDRESS") + ":" + os.Getenv("PORT")

	// Fiber instance
	app := fiber.New(fiber.Config{
		AppName:        "A Simple Api Go Fiber",
		ServerHeader:   "Fiber",
		RequestMethods: []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return ctx.Status(code).JSON(fiber.Map{
				"status": "[router] -> " + strconv.Itoa(code) + " " + err.Error(),
			})
		},
	})
	// Setup App Routes
	router.CreateRoutes(app)

	// Start server
	log.Fatal(app.Listen(listeningAddress))
}
