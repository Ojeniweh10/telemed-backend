package routes

import (
	"telemed/controllers"
	"telemed/responses"

	"github.com/gofiber/fiber/v2"
)

var adminController controllers.AdminController

const (
	Admin = "admin"
)

func AdminRoutes(app *fiber.App) {
	api := app.Group("/admin")
	api.Post("/Login", roleMiddleware(Admin), adminController.Login)
	api.Post("/otp", roleMiddleware(Admin), adminController.VerifyOTP)
	api.Post("/forgot-password", roleMiddleware(Admin), adminController.ForgotPassword)
	api.Post("/verify-forgot-password-otp", roleMiddleware(Admin), adminController.VerifyPwdOTP)
	api.Post("/reset-password", roleMiddleware(Admin), adminController.ResetPassword)
	//api.Get("/get-all-users", roleMiddleware(Admin),middleware.JWTProtected(), adminController.GetAllUsers)
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
