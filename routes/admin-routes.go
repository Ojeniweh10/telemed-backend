package routes

import (
	"telemed/controllers"
	"telemed/middleware"
	"telemed/responses"

	"github.com/gofiber/fiber/v2"
)

var adminController controllers.AdminController

const (
	Admin   = "admin"
	God_eye = "god_eye"
)

func AdminRoutes(app *fiber.App) {
	api := app.Group("/admin")
	api.Post("/Login", roleMiddleware(Admin, God_eye), adminController.Login)
	api.Post("/otp", roleMiddleware(Admin, God_eye), adminController.VerifyOTP)
	api.Post("/forgot-password", roleMiddleware(Admin, God_eye), adminController.ForgotPassword)
	api.Post("/verify-forgot-password-otp", roleMiddleware(Admin, God_eye), adminController.VerifyPwdOTP)
	api.Post("/reset-password", roleMiddleware(Admin, God_eye), adminController.ResetPassword)
	//dashboards
	api.Get("/dashboard/summary", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchDashboardSummary)
	api.Get("/analytics", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAnalytics)
	//appointments
	api.Get("/appointments", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAppointments)
	api.Post("/appointments/:id", roleMiddleware(Admin, God_eye), middleware.JWTProtected(), adminController.FetchAppointmentByID)

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
