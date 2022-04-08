package handlers

import (
	"Assignment1/rest/schema"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"time"
)

type JwtConf struct {
	jwt.StandardClaims
	Data interface{} `json:"data"`
	Time time.Time
}

func JwtEncrypt(c *fiber.Ctx) error {
	data := schema.InputData{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	now := time.Now()

	jwtConf := &JwtConf{
		Data: data.Data,
		Time: now,
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256,jwtConf)
	token,err := rawToken.SignedString([]byte(data.Key))
	if err != nil {
		return err
	}

	return c.JSON(schema.RespondToken{Token: token})
}

func JwtDecrypt(c *fiber.Ctx) error {
	data := schema.InputToken{}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	claims := &JwtConf{}

	_,err := jwt.ParseWithClaims(data.Token,claims, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil,fmt.Errorf("unexpected signing method")
		}

		return []byte(data.Key),nil
	})

	if err != nil {
		return err
	}
	return c.JSON(schema.RespondData{Body: claims.Data,Time: claims.Time})
}