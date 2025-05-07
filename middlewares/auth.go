package middlewares

import (
	
	"github.com/kgermando/ipos-stock-api/utils"
	
	"github.com/gofiber/fiber"
)

func IsAuthenticated(c *fiber.Ctx) error {
 
	token := c.Query("token")


	if _, err := utils.VerifyJwt(token); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
 c.Next()
 return nil
}
