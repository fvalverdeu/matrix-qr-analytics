package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"qr-service-go/configs"
	"qr-service-go/internal/clients"
	"qr-service-go/internal/handlers"
	"qr-service-go/internal/services"
)

func main() {
	cfg := configs.Load()

	validator := services.NewMatrixValidator()
	statisticsClient := clients.NewHTTPStatisticsClient(cfg.StatisticsServiceURL, cfg.StatisticsTimeout)
	qrService := services.NewQRService(validator, statisticsClient)
	qrHandler := handlers.NewQRHandler(qrService)

	app := fiber.New(fiber.Config{
		AppName: "qr-service-go",
	})

	app.Use(recover.New())

	handlers.RegisterRoutes(app, qrHandler)

	log.Printf("QR service listening on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
