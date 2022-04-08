package main

import (
	"Assignment4/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

//Track progress by write the byte manually through the compressor

func main() {

	engine := html.New("./html",".html")

	app := fiber.New(fiber.Config{
		BodyLimit: 1024*1024*1024,
		Views: engine,
	})

	if err := os.Chdir("./static"); err != nil {
		log.Fatal(err.Error())
	}

	app.Use(recover.New())
	
	app.Get("/upload", func(ctx *fiber.Ctx) error {
		return ctx.Render("main",fiber.Map{})
	})

	app.Post("/upload", func(ctx *fiber.Ctx) error {
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

		f,err := os.Create(file.FileName)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		_, err = io.Copy(f,file.File)
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		if err := f.Close(); err != nil {
			return err
		}

		return ctx.Render("download",fiber.Map{
			"Link":fmt.Sprintf("/download/%s",url.PathEscape(file.FileName)),
		})
	})

	if err := app.Listen(":8001"); err !=nil {
		log.Fatal(err.Error())
	}
}