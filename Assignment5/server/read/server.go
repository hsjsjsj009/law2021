package write

import (
	"Assignment5/server/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

func Init(db *gorm.DB) (app *fiber.App) {
	app = fiber.New()

	app.Get("/read/:npm/:trx?", func(ctx *fiber.Ctx) error {
		user := &model.Temp{}
		err := db.Where("npm = ?", ctx.Params("npm")).First(user).Error
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"status":"Not Found",
			})
		}

		return ctx.JSON(fiber.Map{
			"status":"OK",
			"npm":user.Npm,
			"nama": user.Nama,
		})
	})

	return app
}
