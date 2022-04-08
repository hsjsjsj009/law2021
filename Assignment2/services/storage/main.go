package main

import (
	"Assignment2/api"
	"Assignment2/middleware"
	"Assignment2/util"
	"fmt"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"time"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
	})

	workDir := os.Getenv("WORK_DIR")
	if workDir == "" {
		workDir = "./services/storage/storage_folder"
	}

	if err := os.Chdir(workDir);err != nil {
		log.Fatal(err.Error())
	}

	app.Use(recover2.New())
	app.Use(middleware.GoogleAuth)
	app.Static("/assets","./")

	app.Post("/", func(c *fiber.Ctx) error {
		var (
			root = "."
			userData = c.Locals("user_data").(*api.GoogleData)
		)
		formData,err := c.MultipartForm()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		if val,ok := formData.Value["root"];ok {
			root = val[0]
		}

		files,err := util.MultipartFileConverter(formData.File)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		if root != "." {
			root = fmt.Sprintf("%s/%s",userData.ID,root)
		}else {
			root = userData.ID
		}

		if _,err := os.Stat(root); os.IsNotExist(err) {
			err = os.MkdirAll(root,os.ModePerm)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":err.Error(),
				})
			}
		}

		mapFileUrl := map[string]string{}

		for _,file := range files["file"] {
			key := fmt.Sprintf("%d-%s",time.Now().Unix(),file.FileName)
			if root != "." {
				key = fmt.Sprintf("%s/%s",root,key)
			}
			fileWriter,err := os.Create(key)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":err.Error(),
				})
			}
			_,err = io.Copy(fileWriter,file.File)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":err.Error(),
				})
			}

			fileWriter.Close()
			url,err := url2.Parse(key)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message":err.Error(),
				})
			}
			mapFileUrl[file.FileName] = fmt.Sprintf("http://%s/assets/%s",c.Hostname(),url.EscapedPath())
		}

		return c.JSON(mapFileUrl)
	})

	log.Fatal(app.Listen(":8081"))

}