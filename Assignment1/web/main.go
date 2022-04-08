package main

import (
	restSchema "Assignment1/rest/schema"
	serverSchema "Assignment1/web/schema"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	engine:= html.New("./templates",".html")

	restHost := os.Getenv("REST_SERVER")
	if restHost == "" {
		log.Fatal("REST_SERVER required")
	}

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("generate",fiber.Map{
			"Post":false,
		})
	})

	app.Get("/reverse", func(ctx *fiber.Ctx) error {
		return ctx.Render("reverse",fiber.Map{
			"Post":false,
		})
	})

	app.Post("/", func(ctx *fiber.Ctx) error {
		b := serverSchema.InputData{}
		if err := ctx.BodyParser(&b);err != nil {
			return err
		}

		data,err := json.Marshal(b)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/jwt",restHost)
		res,err := http.Post(url,"application/json",bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		resData := &restSchema.RespondToken{}

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(resData)
		if err != nil {
			return err
		}


		return ctx.Render("generate",fiber.Map{
			"Post":true,
			"Token":resData.Token,
		})
	})

	app.Post("/reverse", func(ctx *fiber.Ctx) error {
		b := restSchema.InputToken{}
		if err := ctx.BodyParser(&b);err != nil {
			return err
		}

		data,err := json.Marshal(b)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/jwt-decrypt",restHost)
		res,err := http.Post(url,"application/json",bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		resData := &restSchema.RespondData{}

		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(resData)
		if err != nil {
			return err
		}


		return ctx.Render("reverse",fiber.Map{
			"Post":true,
			"Data":resData.Body,
			"Time":resData.Time.Format(time.RFC822),
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s",port)))
}
