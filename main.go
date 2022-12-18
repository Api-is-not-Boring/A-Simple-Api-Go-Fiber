package main

import (
	"A-Simple-Api-Go-Fiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
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
	})
	// Setup App Routes
	router.CreateRoutes(app)

	// Start server
	log.Fatal(app.Listen(listeningAddress))
}
