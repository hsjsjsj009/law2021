package main

import (
	redis2 "MiniLab1/packages/redis"
	"MiniLab1/packages/str"
	"github.com/go-redis/redis/v7"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	ClientID string	`json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
}

type UserDetailData struct {
	User
	FullName string
	NPM string
}



func main() {

	user := User{
		Username: "username",
		Password: "password",
		ClientID: "clientID",
		ClientSecret: "clientSecret",
	}

	userFullData := UserDetailData{
		User:user,
		FullName: "Dipta Laksmana Baswara Dwiyantoro",
		NPM: "1806235832",
	}

	userList := []UserDetailData{userFullData}

	redisOption := &redis.Options{
		Addr: "localhost:6379",
		Password: "password",
		DB: 0,
	}

	tokenErrorFunc := func(ctx *fiber.Ctx) error {
		return ctx.Status(401).JSON(fiber.Map{
			"error":"invalid_request",
			"Error_description":"ada kesalahan masbro!",
		})
	}

	redisClient := redis2.Client{Client: redis.NewClient(redisOption)}

	app := fiber.New()

	app.Post("/oauth/token", func(ctx *fiber.Ctx) error {
		data := struct {
			User
			GrantType string `json:"grant_type" form:"grant_type"`
		}{}
		contentType := string(ctx.Request().Header.ContentType())
		if contentType != fiber.MIMEApplicationForm {
			return tokenErrorFunc(ctx)
		}
		err := ctx.BodyParser(&data)
		if err != nil {
			return tokenErrorFunc(ctx)
		}

		var userData UserDetailData
		wrongGrantType := data.GrantType != "password"
		userExists := false
		for _,u := range userList {
			if u.Username == data.Username {
				userData = u
				userExists = true
				break
			}
		}
		if wrongGrantType || !userExists {
			return tokenErrorFunc(ctx)
		}

		accessToken := str.RandStringBytesMaskImprSrc(40)
		refreshToken := str.RandStringBytesMaskImprSrc(40)
		if err := redisClient.StoreToRedisWithExpired(
			accessToken,
			fiber.Map{
				"username":userData.Username,
				"refresh_token":refreshToken,
			},
			"5m");err != nil {
			return tokenErrorFunc(ctx)
		}
		return ctx.JSON(fiber.Map{
			"access_token":accessToken,
			"expires_in":300,
			"token_type":"Bearer",
			"scope":nil,
			"refresh_token":refreshToken,
		})
	})

	resourceError := func(c *fiber.Ctx) error {
		return c.Status(401).JSON(fiber.Map{
			"error":"invalid_token",
			"error_description":"Token Salah masbro",
		})
	}

	app.Post("/oauth/resource", func(ctx *fiber.Ctx) (err error) {
		defer func() {
			if errRec := recover(); errRec != nil {
				err = resourceError(ctx)
			}
		}()
		headerBearer := ctx.Get("Authorization","")
		if headerBearer == "" || !strings.Contains(headerBearer,"Bearer") {
			return resourceError(ctx)
		}
		contentType := string(ctx.Request().Header.ContentType())
		if contentType != fiber.MIMEApplicationForm {
			return resourceError(ctx)
		}
		splitHeader := strings.Split(headerBearer," ")
		token := splitHeader[1]
		if token == "" {
			return resourceError(ctx)
		}

		mapData := fiber.Map{}

		if err := redisClient.GetFromRedis(token,&mapData);err != nil {
			return resourceError(ctx)
		}

		var userData UserDetailData
		for _,u := range userList {
			if u.Username == mapData["username"] {
				userData = u
				break
			}
		}

		return ctx.JSON(fiber.Map{
			"access_token":token,
			"client_id":userData.ClientID,
			"user_id":userData.Username,
			"full_name":userData.FullName,
			"npm":userData.NPM,
			"expires":nil,
			"refresh_token":mapData["refresh_token"],
		})
	})

	log.Fatal(app.Listen(":8080"))
}
