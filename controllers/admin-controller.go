package controllers

import (
	"telemed/models"
	"telemed/responses"
	"telemed/servers"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct{}

var adminServer servers.AdminServer

func (AdminController) Login(c *fiber.Ctx) error {
	var payload models.Adminlogin
	//parse data from request
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	//vet if data exists
	if payload.Email == "" || payload.Password == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	//pass data to servers
	res, err := adminServer.Login(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (AdminController) VerifyOTP(c *fiber.Ctx) error {
	var payload models.OTPVerify
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.OTP == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.VerifyOTP(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_VERIFIED, res, 200)
}

func (AdminController) ForgotPassword(c *fiber.Ctx) error {
	var payload models.ForgotPassword
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.ForgotPassword(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_SENT, res, 200)
}

func (AdminController) VerifyPwdOTP(c *fiber.Ctx) error {
	var payload models.VerifyPwdOTP
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.OTP == "" || payload.Email == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.VerifyPwdOTP(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.OTP_VERIFIED, res, 200)
}

func (AdminController) ResetPassword(c *fiber.Ctx) error {
	var payload models.ResetPassword
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.Email == "" || payload.NewPassword == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.ResetPassword(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.PASSWORD_RESET_SUCCESS, res, 200)
}
func (AdminController) FetchDashboardSummary(c *fiber.Ctx) error {
	res, err := adminServer.GetDashboardSummary()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}
func (AdminController) FetchAnalytics(c *fiber.Ctx) error {
	var data models.AnalyticsReq
	data.Metric = c.Query("metric")
	data.Month = c.Query("month")
	data.Year = c.Query("year")

	if data.Metric == "" || data.Month == "" || data.Year == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	res, err := adminServer.GetAnalytics(data)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)

}

func (AdminController) FetchAppointments(c *fiber.Ctx) error {
	res, err := adminServer.GetAppointments()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchAppointmentByID(c *fiber.Ctx) error {
	var payload models.AppointmentID
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	if payload.ID == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetAppointmentByID(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchDoctorByID(c *fiber.Ctx) error {
	var payload models.Doctorreq
	payload.DoctorTag = c.Params("id")
	if payload.DoctorTag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.GetDoctorByID(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) UpdateAppointmentStatus(c *fiber.Ctx) error {
	var payload models.UpdateAppointmentStatus
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}
	payload.Appointment_id = c.Params("id")
	if payload.Appointment_id == "" || payload.Status == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	res, err := adminServer.UpdateAppointmentStatus(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (a *AdminController) UpdateAppointment(c *fiber.Ctx) error {
	var payload models.RescheduleAppointmentReq
	if err := c.BodyParser(&payload); err != nil {
		return responses.ErrorResponse(c, responses.BAD_DATA, 400)
	}

	payload.Appointment_id = c.Params("id")
	if payload.Appointment_id == "" || payload.NewScheduledAt == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}

	res, err := adminServer.RescheduleAppointment(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}

	return responses.SuccessResponse(c, responses.DATA_UPDATED, res, 200)
}

func (AdminController) FetchDoctors(c *fiber.Ctx) error {
	res, err := adminServer.GetDoctors()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) DeleteDoctor(c *fiber.Ctx) error {
	var payload models.Doctorreq
	payload.DoctorTag = c.Params("id")
	if payload.DoctorTag == "" {
		return responses.ErrorResponse(c, responses.INCOMPLETE_DATA, 400)
	}
	err := adminServer.DeleteDoctor(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_UPDATED, nil, 200)
}

func (AdminController) FetchPatients(c *fiber.Ctx) error {
	res, err := adminServer.GetPatients()
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}

func (AdminController) FetchPatientByUsertag(c *fiber.Ctx) error {
	var payload models.PatientIdReq
	payload.Usertag = c.Params("usertag")

	res, err := adminServer.GetPatientByUsertag(payload)
	if err != nil {
		return responses.ErrorResponse(c, err.Error(), 400)
	}
	return responses.SuccessResponse(c, responses.DATA_FETCHED, res, 200)
}
