package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1024*1024*1024,
	})

	app.Use(recover.New())
	app.Static("/download","./static")

	if err := app.Listen(":8002"); err !=nil {
		log.Fatal(err.Error())
	}
}
