package handlers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"qr-service-go/internal/clients"
	"qr-service-go/internal/models"
	"qr-service-go/internal/services"
)

type QRHandler struct {
	qrService *services.QRService
}

func NewQRHandler(qrService *services.QRService) *QRHandler {
	return &QRHandler{qrService: qrService}
}

func (h *QRHandler) Decompose(c *fiber.Ctx) error {
	var req models.QRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Request body must be valid JSON with a matrix field",
			},
		})
	}

	response, err := h.qrService.Process(req.Matrix)
	if err != nil {
		var validationErr *services.ValidationError
		if errors.As(err, &validationErr) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error: models.APIError{
					Code:    "INVALID_MATRIX",
					Message: validationErr.Message,
				},
			})
		}

		var unavailableErr *clients.UnavailableError
		if errors.As(err, &unavailableErr) {
			return c.Status(fiber.StatusBadGateway).JSON(models.ErrorResponse{
				Error: models.APIError{
					Code:    "STATISTICS_UNAVAILABLE",
					Message: unavailableErr.Message,
				},
			})
		}

		log.Printf("qr_handler: unexpected processing error method=%s path=%s error_type=%T", c.Method(), c.Path(), err)

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error: models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "An unexpected error occurred while processing the matrix",
			},
		})
	}

	return c.JSON(response)
}
