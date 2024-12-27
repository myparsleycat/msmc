package routes

import (
	"msmc/src/server/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UtilsRoute(router fiber.Router, db *gorm.DB) {
	utils := router.Group("/utils")

	utils.Get("/reload", func(c *fiber.Ctx) error {
		return handlers.UtilsReloadHandler(c, db)
	})

	utils.Get("/test", func(c *fiber.Ctx) error {
		return handlers.UtilsTestHandler(c, db)
	})
}
