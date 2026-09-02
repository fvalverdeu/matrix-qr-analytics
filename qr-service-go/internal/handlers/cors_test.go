package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"qr-service-go/internal/models"
	"qr-service-go/internal/services"
)

func TestCORS_AllowedOriginReceivesHeaders(t *testing.T) {
	fakeClient := &fakeStatisticsClient{
		result: models.Statistics{
			Max:               6,
			Min:               1,
			Average:           3.5,
			Sum:               14,
			HasDiagonalMatrix: false,
		},
	}

	app := buildTestAppWithCORS(fakeClient, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/qr", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected Access-Control-Allow-Origin: got %q", got)
	}

	if got := resp.Header.Get("Access-Control-Allow-Methods"); got != "GET,POST,OPTIONS" {
		t.Fatalf("unexpected Access-Control-Allow-Methods: got %q", got)
	}

	if got := resp.Header.Get("Access-Control-Allow-Headers"); got != "Content-Type" {
		t.Fatalf("unexpected Access-Control-Allow-Headers: got %q", got)
	}
}

func TestCORS_UnlistedOriginDoesNotReceiveAllowOriginHeader(t *testing.T) {
	fakeClient := &fakeStatisticsClient{}
	app := buildTestAppWithCORS(fakeClient, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/qr", nil)
	req.Header.Set("Origin", "http://localhost:3001")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin for unlisted origin, got %q", got)
	}
}

func TestCORS_ExistingRouteBehaviorIsPreserved(t *testing.T) {
	fakeClient := &fakeStatisticsClient{
		result: models.Statistics{
			Max:               4,
			Min:               0,
			Average:           1.375,
			Sum:               11,
			HasDiagonalMatrix: true,
		},
	}

	app := buildTestAppWithCORS(fakeClient, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(`{"matrix":[[1,2,3],[4,5,6]]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:5173")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected Access-Control-Allow-Origin: got %q", got)
	}

	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client call count 1, got %d", fakeClient.callCount)
	}
}

func buildTestAppWithCORS(statisticsClient services.StatisticsClient, allowedOrigins string) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	validator := services.NewMatrixValidator()
	qrService := services.NewQRService(validator, statisticsClient)
	qrHandler := NewQRHandler(qrService)
	RegisterRoutes(app, qrHandler)

	return app
}
