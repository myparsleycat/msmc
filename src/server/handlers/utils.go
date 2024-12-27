package handlers

import (
	"msmc/src/browser"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UtilsReloadHandler(c *fiber.Ctx, db *gorm.DB) error {
	err := browser.TriggerPageReload()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "페이지 새로고침 성공",
	})
}
