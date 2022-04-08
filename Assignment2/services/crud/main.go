package main

import (
	"Assignment2/api"
	"Assignment2/middleware"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserData struct {
	Password string `json:"password"`
	Address string `json:"address"`
	BirthDate time.Time `json:"birth_date"`
}

type UserCompleteData struct {
	*User
	Files []*File `json:"files"`
}

type User struct {
	UserData
	Email string `json:"email"`
	ID int `json:"id"`
}

type FileData struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Description string `json:"description"`
}

type File struct {
	FileData
	UserEmail string `json:"user_email"`
	ID        int `json:"id"`
}

func main() {

	userCount := 0
	fileCount := 0
	var listUsers []*User
	var listFiles []*File

	app := fiber.New()

	app.Use(recover2.New())
	app.Use(middleware.GoogleAuth)

	app.Get("/file/:id", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)
		id,err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"id must be integer",
			})
		}

		var fileData *File
		for _, file := range listFiles {
			if file.ID == id {
				fileData = file
				break
			}
		}

		if fileData == nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"file not found",
			})
		}

		var userData *User
		for _,user := range listUsers {
			if user.Email == fileData.UserEmail {
				userData = user
				break
			}
		}

		if userData == nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"user not found",
			})
		}

		if userData.Email != authData.Email {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"message":"not allowed",
			})
		}

		return c.JSON(fileData)
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)

		var userData *User
		for _,user := range listUsers {
			if user.Email == authData.Email {
				userData = user
				break
			}
		}

		if userData == nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"user not found",
			})
		}

		var userFiles []*File
		for _,file := range listFiles {
			if file.UserEmail == userData.Email {
				userFiles = append(userFiles,file)
			}
		}

		data := UserCompleteData{
			User: userData,
			Files: userFiles,
		}

		return c.JSON(data)

	})

	app.Post("/user", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)
		data := UserData{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		exist := false
		for _,user := range listUsers {
			if user.Email == authData.Email {
				exist = true
				break
			}
		}

		if exist {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"already registered",
			})
		}

		userCount += 1

		listUsers = append(listUsers,&User{ID: userCount,Email:authData.Email,UserData:data})

		return c.JSON(fiber.Map{
			"message":"success",
		})
	})

	app.Post("/file", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)

		data := FileData{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		fileCount += 1
		id := fileCount
		listFiles = append(listFiles,&File{ID: id, UserEmail:authData.Email,FileData:data})

		return c.JSON(fiber.Map{
			"message":"success",
			"id":fileCount,
		})
	})

	app.Put("/user", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)

		data := UserData{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		exists := false
		for _,user := range listUsers {
			if user.Email == authData.Email {
				user.UserData = data
				exists = true
				break
			}
		}

		if !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"user not found",
			})
		}

		return c.JSON(fiber.Map{
			"message":"success",
		})
	})

	app.Put("/file/:id", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)
		id,err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"id must be integer",
			})
		}

		data := FileData{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message":err.Error(),
			})
		}

		exists := false
		for _, file := range listFiles {
			if file.ID == id {
				if file.UserEmail != authData.Email {
					return c.Status(http.StatusForbidden).JSON(fiber.Map{
						"message":"not allowed",
					})
				}
				file.FileData = data
				exists = true
				break
			}
		}

		if !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"file not found",
			})
		}

		return c.JSON(fiber.Map{
			"message":"success",
		})
	})

	app.Delete("/file/:id", func(c *fiber.Ctx) error {
		authData := c.Locals("user_data").(*api.GoogleData)
		id,err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"id must be integer",
			})
		}

		exists := false
		var newFilesList []*File
		for _, file := range listFiles {
			if file.ID != id {
				newFilesList = append(newFilesList,file)
			}else if file.UserEmail != authData.Email{
					return c.Status(http.StatusForbidden).JSON(fiber.Map{
						"message":"not allowed",
				})
			}
			exists = true
		}

		if !exists {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"message":"file not found",
			})
		}

		listFiles = newFilesList

		return c.JSON(fiber.Map{
			"message":"success",
		})
	})

	log.Fatal(app.Listen(":8082"))



}