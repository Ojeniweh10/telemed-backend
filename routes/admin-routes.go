package routes

import (
	"telemed/responses"

	"github.com/gofiber/fiber/v2"
)

const (
	Admin = "admin"
)

func AdminRoutes(app *fiber.App) {
	api := app.Group("/admin")
	api.Get("/", roleMiddleware(Admin))
}

func roleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminRole := c.Get("role")
		if !contains(allowedRoles, adminRole) {
			return responses.ErrorResponse(c, "Unauthorized access", fiber.StatusForbidden)
		}
		return c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
