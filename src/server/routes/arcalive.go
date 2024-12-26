package routes

import (
	"github.com/gofiber/fiber/v2"
	"msmc/src/server/handlers"
)

func SetupArcaliveRoutes(router fiber.Router) {
	arcalive := router.Group("/arcalive")

	arcalive.Get("/:username", handlers.GetUserState)
}
