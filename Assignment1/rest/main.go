package main

import (
	"Assignment1/rest/handlers"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"
)

func main() {
	loc,err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Local = loc

	app := fiber.New()

	app.Post("/jwt", handlers.JwtEncrypt)
	app.Post("/jwt-decrypt", handlers.JwtDecrypt)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s",port)))
}
