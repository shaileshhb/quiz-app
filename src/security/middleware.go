package security

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MandatoryAuthMiddleware will check that authorization cookie is valid.
func MandatoryAuthMiddleware(c *fiber.Ctx) error {
	authorizationTypeBearer := "bearer"

	authHeader := c.Get("authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header provided",
		})
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != authorizationTypeBearer {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("unsupported authorization type %s", authorizationType),
		})
	}

	user, err := ValidateJWT(fields[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	c.Locals("user", user)
	return c.Next()
}
