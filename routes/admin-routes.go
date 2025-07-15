package routes

import (
	"telemed/controllers"
	"telemed/middleware"
	"telemed/responses"

	"github.com/gofiber/fiber/v2"
)

var adminController controllers.AdminController

const (
	Admin = "admin"
)

func AdminRoutes(app *fiber.App) {
	api := app.Group("/admin")
	api.Use(middleware.JWTProtected())
	api.Post("/Login", roleMiddleware(Admin), adminController.Login)
	api.Post("/otp", roleMiddleware(Admin), adminController.VerifyOTP)
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
