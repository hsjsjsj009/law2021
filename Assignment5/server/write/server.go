package write

import (
	"Assignment5/server/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) (app *fiber.App) {
	app = fiber.New()

	app.Post("/update", func(ctx *fiber.Ctx) error {
		user := &model.Temp{}
		err := ctx.BodyParser(user)
		if err != nil {
			return err
		}

		db.FirstOrCreate(user, *user)
		return ctx.JSON(fiber.Map{
			"status":"OK",
		})
	})

	return app
}
