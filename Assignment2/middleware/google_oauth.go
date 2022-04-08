package middleware

import (
	"Assignment2/api"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func authError(c *fiber.Ctx,msg string) error {
	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"message":msg,
	})
}

func GoogleAuth(c *fiber.Ctx) error {
	var token string
	bearerToken := c.Get("Authorization")
	if bearerToken == "" {
		if c.Method() != http.MethodGet {
			return authError(c,"need authorization header")
		}
		token = c.Query("access_token")
		if token == "" {
			return authError(c,"need authorization header")
		}
	}else {
		bearerSplit := strings.Split(bearerToken," ")
		token = bearerSplit[1]
	}

	googleApi,err := api.NewGoogleAPI(token)
	if err != nil {
		return authError(c,err.Error())
	}
	data,err := googleApi.GetUserData()
	if err != nil {
		return authError(c,err.Error())
	}
	c.Locals("user_data",data)
	return c.Next()
}
