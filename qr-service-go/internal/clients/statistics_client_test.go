package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"qr-service-go/internal/models"
)

func TestHTTPStatisticsClient_CalculateStatistics_SuccessAndRequestContract(t *testing.T) {
	inputQ := [][]float64{{1, 0}, {0, 1}}
	inputR := [][]float64{{2, 3, 4}, {0, 5, 6}}

	type capturedRequest struct {
		method      string
		path        string
		contentType string
		body        models.StatisticsRequest
		decodeErr   error
	}

	capturedCh := make(chan capturedRequest, 1)

	expected := models.Statistics{
		Max:               10,
		Min:               -2,
		Average:           3.5,
		Sum:               21,
		HasDiagonalMatrix: false,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured := capturedRequest{
			method:      r.Method,
			path:        r.URL.Path,
			contentType: r.Header.Get("Content-Type"),
		}
		var req models.StatisticsRequest
		captured.decodeErr = json.NewDecoder(r.Body).Decode(&req)
		captured.body = req
		capturedCh <- captured

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":10,"min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":false}`))
	}))
	defer server.Close()

	client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)

	stats, err := client.CalculateStatistics(inputQ, inputR)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if stats != expected {
		t.Fatalf("unexpected statistics: got %+v, want %+v", stats, expected)
	}

	var captured capturedRequest
	select {
	case captured = <-capturedCh:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for captured request")
	}

	if captured.method != http.MethodPost {
		t.Fatalf("unexpected method: got %s, want %s", captured.method, http.MethodPost)
	}
	if captured.path != "/api/v1/statistics" {
		t.Fatalf("unexpected path: got %s, want /api/v1/statistics", captured.path)
	}
	if captured.contentType != "application/json" {
		t.Fatalf("unexpected Content-Type: got %q, want %q", captured.contentType, "application/json")
	}
	if captured.decodeErr != nil {
		t.Fatalf("failed to decode request body: %v", captured.decodeErr)
	}
	if !reflect.DeepEqual(captured.body.Q, inputQ) {
		t.Fatalf("unexpected request Q: got %v, want %v", captured.body.Q, inputQ)
	}
	if !reflect.DeepEqual(captured.body.R, inputR) {
		t.Fatalf("unexpected request R: got %v, want %v", captured.body.R, inputR)
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_Non200ReturnsUnavailableError(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{name: "status_400", statusCode: http.StatusBadRequest},
		{name: "status_500", statusCode: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(`{"error":"downstream"}`))
			}))
			defer server.Close()

			client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)
			stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			assertUnavailableError(t, err)
			if stats != (models.Statistics{}) {
				t.Fatalf("expected zero-value statistics on error, got %+v", stats)
			}
		})
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_InvalidSuccessfulResponseBodyReturnsUnavailableError(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{name: "malformed_json", body: `{"max":10,`},
		{name: "empty_body", body: ``},
		{name: "wrong_field_type", body: `{"max":"not-a-number","min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":false}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)
			stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			assertUnavailableError(t, err)
			if stats != (models.Statistics{}) {
				t.Fatalf("expected zero-value statistics on error, got %+v", stats)
			}
		})
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_TransportFailureReturnsUnavailableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close()

	client := NewHTTPStatisticsClient(url, 500*time.Millisecond)
	stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertUnavailableError(t, err)
	if stats != (models.Statistics{}) {
		t.Fatalf("expected zero-value statistics on error, got %+v", stats)
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_TimeoutReturnsUnavailableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":10,"min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":false}`))
	}))
	defer server.Close()

	client := NewHTTPStatisticsClient(server.URL, 50*time.Millisecond)
	stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	assertUnavailableError(t, err)
	if stats != (models.Statistics{}) {
		t.Fatalf("expected zero-value statistics on error, got %+v", stats)
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_ExtraResponseFieldsAreIgnored(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":10,"min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":false,"futureField":"ignored"}`))
	}))
	defer server.Close()

	client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)
	stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expected := models.Statistics{Max: 10, Min: -2, Average: 3.5, Sum: 21, HasDiagonalMatrix: false}
	if stats != expected {
		t.Fatalf("unexpected statistics: got %+v, want %+v", stats, expected)
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_ZeroValuesPresentResponseIsValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":0,"min":0,"average":0,"sum":0,"hasDiagonalMatrix":false}`))
	}))
	defer server.Close()

	client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)
	stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expected := models.Statistics{
		Max:               0,
		Min:               0,
		Average:           0,
		Sum:               0,
		HasDiagonalMatrix: false,
	}
	if stats != expected {
		t.Fatalf("unexpected statistics: got %+v, want %+v", stats, expected)
	}
}

func TestHTTPStatisticsClient_CalculateStatistics_IncompleteSuccessfulResponseMustReturnError(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{name: "only_max", body: `{"max":10}`},
		{name: "empty_object", body: `{}`},
		{name: "missing_has_diagonal_matrix", body: `{"max":10,"min":-2,"average":3.5,"sum":21}`},
		{name: "null_has_diagonal_matrix", body: `{"max":10,"min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":null}`},
		{name: "null_numeric_field", body: `{"max":null,"min":-2,"average":3.5,"sum":21,"hasDiagonalMatrix":false}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := NewHTTPStatisticsClient(server.URL, 500*time.Millisecond)
			stats, err := client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})

			if err == nil {
				t.Fatalf("expected error for incomplete/null statistics response, got nil and stats=%+v", stats)
			}
			assertUnavailableError(t, err)
			if stats != (models.Statistics{}) {
				t.Fatalf("expected zero-value statistics on error, got %+v", stats)
			}
		})
	}
}

func assertUnavailableError(t *testing.T, err error) {
	t.Helper()

	var unavailableErr *UnavailableError
	if !errors.As(err, &unavailableErr) {
		t.Fatalf("expected *UnavailableError, got %T (%v)", err, err)
	}
}

func TestNewHTTPStatisticsClientWithAuth_NoneModePreservesBehavior(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Fatalf("did not expect Authorization header in none mode")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":1,"min":1,"average":1,"sum":1,"hasDiagonalMatrix":true}`))
	}))
	defer server.Close()

	client, err := NewHTTPStatisticsClientWithAuth(context.Background(), server.URL, 750*time.Millisecond, StatisticsAuthModeNone)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if client.httpClient.Timeout != 750*time.Millisecond {
		t.Fatalf("expected timeout %v, got %v", 750*time.Millisecond, client.httpClient.Timeout)
	}

	_, err = client.CalculateStatistics([][]float64{{1}}, [][]float64{{1}})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNewHTTPStatisticsClientWithAuth_GoogleIDTokenUsesFactoryAndPreservesTimeout(t *testing.T) {
	originalFactory := newIDTokenHTTPClient
	t.Cleanup(func() {
		newIDTokenHTTPClient = originalFactory
	})

	called := false
	gotAudience := ""
	newIDTokenHTTPClient = func(ctx context.Context, audience string) (*http.Client, error) {
		called = true
		gotAudience = audience
		return &http.Client{}, nil
	}

	client, err := NewHTTPStatisticsClientWithAuth(context.Background(), "https://example.run.app/", 900*time.Millisecond, StatisticsAuthModeGoogleIDToken)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !called {
		t.Fatal("expected Google ID token client factory to be called")
	}
	if gotAudience != "https://example.run.app" {
		t.Fatalf("expected audience to be base URL without trailing slash, got %q", gotAudience)
	}
	if client.httpClient.Timeout != 900*time.Millisecond {
		t.Fatalf("expected timeout %v, got %v", 900*time.Millisecond, client.httpClient.Timeout)
	}
}

func TestNewHTTPStatisticsClientWithAuth_UnsupportedModeReturnsError(t *testing.T) {
	_, err := NewHTTPStatisticsClientWithAuth(context.Background(), "https://example.run.app", 500*time.Millisecond, "unsupported")
	if err == nil {
		t.Fatal("expected error for unsupported auth mode")
	}
}
