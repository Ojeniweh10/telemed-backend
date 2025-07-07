package responses

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type squadResponse struct {
	Response_code         int    `json:"response_code"`
	Transaction_reference string `json:"transaction_reference"`
	Response_description  string `json:"response_description"`
}

func ErrorResponse(c *fiber.Ctx, message string, statusCode int) error {
	res := Response{
		Success: false,
		Message: message,
	}
	return c.Status(statusCode).JSON(res)
}

func SuccessResponse(c *fiber.Ctx, message string, data any, statusCode int) error {
	res := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(res)
}

func SquadResponse(c *fiber.Ctx, reference, description string, statusCode int) error {
	res := squadResponse{
		Response_code:         statusCode,
		Transaction_reference: reference,
		Response_description:  description,
	}
	return c.Status(statusCode).JSON(res)
}

const (
	UNAUTHORIZED_ACCESS = "unauthorized access"
)
