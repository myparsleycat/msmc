package routes

import (
	"msmc/src/server/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupArcaliveRoutes(router fiber.Router, db *gorm.DB) {
	arcalive := router.Group("/arcalive")

	// arcalive.Get("/audit", handlers.GetAuditsHandler)
	arcalive.Get("/:username", func(c *fiber.Ctx) error {
		return handlers.GetUserState(c, db)
	})
}
