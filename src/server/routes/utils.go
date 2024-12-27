package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UtilsRoute(router fiber.Router, db *gorm.DB) {
	utils := router.Group("/utils")

	utils.Get("/reload")
}
