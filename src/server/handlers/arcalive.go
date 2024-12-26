package handlers

import (
	"msmc/src/database"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

type GetUserStateHttpResponse struct {
	Posts   int64 `json:"posts"`
	Reviews int64 `json:"reviews"`
}

type GetUserStateQueryResult struct {
	TotalCount  int64
	ReviewCount int64
}

func GetUserState(c *fiber.Ctx) error {
	db := database.GetDB()

	encodedUsername := c.Params("username")
	username, err := url.QueryUnescape(encodedUsername)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid username encoding",
		})
	}

	var result GetUserStateQueryResult
	db.Raw(`
		SELECT 
			COUNT(*) as total_count,
			SUM(CASE WHEN category = '모드리뷰' THEN 1 ELSE 0 END) as review_count 
		FROM Post 
		WHERE LOWER(author) = LOWER(?)
	`, username).Scan(&result)

	response := GetUserStateHttpResponse{
		Posts:   result.TotalCount,
		Reviews: result.ReviewCount,
	}

	return c.JSON(response)
}
