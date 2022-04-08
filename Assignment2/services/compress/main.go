package main

import (
	"Assignment2/api"
	appMiddleware "Assignment2/middleware"
	"Assignment2/util"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 40 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(appMiddleware.GoogleAuth)

	app.Post("/", func(ctx *fiber.Ctx) error {
		userData := ctx.Locals("user_data").(*api.GoogleData)

		form,err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		files,err := util.MultipartFileConverter(form.File)

		if len(files["file"]) > 1 {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"only process one file at a time",
			})
		}

		file := files["file"][0]

		fileName := strings.Split(file.FileName,".")

		fileBytes,err := ioutil.ReadAll(file.File)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		ctx.Set("Content-Type","application/gzip")

		writerBuffer := new(bytes.Buffer)


		gz,err := gzip.NewWriterLevel(writerBuffer,gzip.BestCompression)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		gz.Name = fmt.Sprintf("%s-%s",userData.Email,file.FileName)

		_,err = gz.Write(fileBytes)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		_ = gz.Close()

		ctx.Set("Content-Disposition",fmt.Sprintf(`filename="%s-%s.gz"`,userData.Email,fileName[0]))

		return ctx.SendStream(writerBuffer)
	})

	log.Fatal(app.Listen(":8080"))
}