package handlers

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, qrHandler *QRHandler) {
	app.Get("/health", healthHandler)
	app.Post("/api/v1/qr", qrHandler.Decompose)
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "qr-service-go",
	})
}
