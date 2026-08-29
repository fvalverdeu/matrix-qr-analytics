package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"qr-service-go/internal/clients"
	"qr-service-go/internal/models"
	"qr-service-go/internal/services"
)

type fakeStatisticsClient struct {
	callCount int
	capturedQ [][]float64
	capturedR [][]float64
	result    models.Statistics
	err       error
}

func (f *fakeStatisticsClient) CalculateStatistics(q, r [][]float64) (models.Statistics, error) {
	f.callCount++
	f.capturedQ = q
	f.capturedR = r
	if f.err != nil {
		return models.Statistics{}, f.err
	}
	return f.result, nil
}

func TestQRHandler_Decompose_Success(t *testing.T) {
	fakeClient := &fakeStatisticsClient{
		result: models.Statistics{
			Max:               10,
			Min:               -2,
			Average:           3.5,
			Sum:               21,
			HasDiagonalMatrix: false,
		},
	}
	app := buildTestApp(fakeClient)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(`{"matrix":[[1,2],[3,4]]}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var body models.QRResponse
	decodeResponseJSON(t, resp, &body)

	if body.Statistics != fakeClient.result {
		t.Fatalf("unexpected statistics: got %+v, want %+v", body.Statistics, fakeClient.result)
	}
	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client call count 1, got %d", fakeClient.callCount)
	}
	if !reflect.DeepEqual(fakeClient.capturedQ, body.Q) {
		t.Fatalf("captured Q differs from response Q")
	}
	if !reflect.DeepEqual(fakeClient.capturedR, body.R) {
		t.Fatalf("captured R differs from response R")
	}
	if len(body.Q) == 0 || len(body.Q[0]) == 0 {
		t.Fatalf("response Q must be non-empty")
	}
	if len(body.R) == 0 || len(body.R[0]) == 0 {
		t.Fatalf("response R must be non-empty")
	}
}

func TestQRHandler_Decompose_InvalidRequestBody(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{name: "malformed_json", body: `{"matrix":[[1,2]`},
		{name: "wrong_matrix_type", body: `{"matrix":"not-a-matrix"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeClient := &fakeStatisticsClient{}
			app := buildTestApp(fakeClient)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test returned error: %v", err)
			}

			if resp.StatusCode != http.StatusBadRequest {
				t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusBadRequest)
			}

			var body models.ErrorResponse
			decodeResponseJSON(t, resp, &body)

			if body.Error.Code != "INVALID_REQUEST" {
				t.Fatalf("unexpected error code: got %q, want %q", body.Error.Code, "INVALID_REQUEST")
			}
			if body.Error.Message != "Request body must be valid JSON with a matrix field" {
				t.Fatalf("unexpected error message: got %q", body.Error.Message)
			}
			if fakeClient.callCount != 0 {
				t.Fatalf("expected statistics client call count 0, got %d", fakeClient.callCount)
			}
		})
	}
}

func TestQRHandler_Decompose_InvalidMatrix(t *testing.T) {
	fakeClient := &fakeStatisticsClient{}
	app := buildTestApp(fakeClient)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(`{"matrix":[]}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	var body models.ErrorResponse
	decodeResponseJSON(t, resp, &body)

	if body.Error.Code != "INVALID_MATRIX" {
		t.Fatalf("unexpected error code: got %q, want %q", body.Error.Code, "INVALID_MATRIX")
	}
	if body.Error.Message != "Matrix cannot be empty" {
		t.Fatalf("unexpected error message: got %q, want %q", body.Error.Message, "Matrix cannot be empty")
	}
	if fakeClient.callCount != 0 {
		t.Fatalf("expected statistics client call count 0, got %d", fakeClient.callCount)
	}
}

func TestQRHandler_Decompose_StatisticsUnavailable(t *testing.T) {
	fakeClient := &fakeStatisticsClient{
		err: &clients.UnavailableError{Message: "Statistics service is unavailable"},
	}
	app := buildTestApp(fakeClient)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(`{"matrix":[[1,2],[3,4]]}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusBadGateway)
	}

	var body models.ErrorResponse
	decodeResponseJSON(t, resp, &body)

	if body.Error.Code != "STATISTICS_UNAVAILABLE" {
		t.Fatalf("unexpected error code: got %q, want %q", body.Error.Code, "STATISTICS_UNAVAILABLE")
	}
	if body.Error.Message != "Statistics service is unavailable" {
		t.Fatalf("unexpected error message: got %q, want %q", body.Error.Message, "Statistics service is unavailable")
	}
	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client call count 1, got %d", fakeClient.callCount)
	}
}

func TestQRHandler_Decompose_InternalError(t *testing.T) {
	fakeClient := &fakeStatisticsClient{err: errors.New("unexpected failure")}
	app := buildTestApp(fakeClient)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/qr", bytes.NewBufferString(`{"matrix":[[1,2],[3,4]]}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test returned error: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	responseBytes, body := decodeErrorResponseJSON(t, resp)

	if body.Error.Code != "INTERNAL_ERROR" {
		t.Fatalf("unexpected error code: got %q, want %q", body.Error.Code, "INTERNAL_ERROR")
	}
	if body.Error.Message != "An unexpected error occurred while processing the matrix" {
		t.Fatalf("unexpected error message: got %q", body.Error.Message)
	}
	if strings.Contains(string(responseBytes), "unexpected failure") {
		t.Fatalf("response body leaks internal error details")
	}
	if fakeClient.callCount != 1 {
		t.Fatalf("expected statistics client call count 1, got %d", fakeClient.callCount)
	}
}

func buildTestApp(statisticsClient services.StatisticsClient) *fiber.App {
	app := fiber.New()
	validator := services.NewMatrixValidator()
	qrService := services.NewQRService(validator, statisticsClient)
	qrHandler := NewQRHandler(qrService)
	RegisterRoutes(app, qrHandler)
	return app
}

func decodeResponseJSON(t *testing.T, resp *http.Response, dst interface{}) {
	t.Helper()
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if err := json.Unmarshal(payload, dst); err != nil {
		t.Fatalf("failed to unmarshal response body: %v; body=%s", err, string(payload))
	}
}

func decodeErrorResponseJSON(t *testing.T, resp *http.Response) ([]byte, models.ErrorResponse) {
	t.Helper()
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var body models.ErrorResponse
	if err := json.Unmarshal(payload, &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v; body=%s", err, string(payload))
	}

	return payload, body
}
