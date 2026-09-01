package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/idtoken"

	"qr-service-go/internal/models"
)

const (
	StatisticsAuthModeNone          = "none"
	StatisticsAuthModeGoogleIDToken = "google-id-token"
)

var newIDTokenHTTPClient = func(ctx context.Context, audience string) (*http.Client, error) {
	return idtoken.NewClient(ctx, audience)
}

type HTTPStatisticsClient struct {
	baseURL    string
	httpClient *http.Client
}

type statisticsResponseDTO struct {
	Max               *float64 `json:"max"`
	Min               *float64 `json:"min"`
	Average           *float64 `json:"average"`
	Sum               *float64 `json:"sum"`
	HasDiagonalMatrix *bool    `json:"hasDiagonalMatrix"`
}

func NewHTTPStatisticsClient(baseURL string, timeout time.Duration) *HTTPStatisticsClient {
	client, _ := NewHTTPStatisticsClientWithAuth(context.Background(), baseURL, timeout, StatisticsAuthModeNone)
	return client
}

func NewHTTPStatisticsClientWithAuth(ctx context.Context, baseURL string, timeout time.Duration, authMode string) (*HTTPStatisticsClient, error) {
	trimmedBaseURL := strings.TrimRight(baseURL, "/")

	httpClient, err := buildHTTPClient(ctx, trimmedBaseURL, authMode)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = timeout

	return &HTTPStatisticsClient{
		baseURL:    trimmedBaseURL,
		httpClient: httpClient,
	}, nil
}

func buildHTTPClient(ctx context.Context, audience string, authMode string) (*http.Client, error) {
	switch authMode {
	case StatisticsAuthModeNone:
		return &http.Client{}, nil
	case StatisticsAuthModeGoogleIDToken:
		client, err := newIDTokenHTTPClient(ctx, audience)
		if err != nil {
			return nil, fmt.Errorf("failed to create Google ID token client for statistics service: %w", err)
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unsupported statistics auth mode %q", authMode)
	}
}

func (c *HTTPStatisticsClient) CalculateStatistics(q, r [][]float64) (models.Statistics, error) {
	if c.baseURL == "" {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	payload, err := json.Marshal(models.StatisticsRequest{Q: q, R: r})
	if err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	url := c.baseURL + "/api/v1/statistics"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	if resp.StatusCode != http.StatusOK {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	var statisticsDTO statisticsResponseDTO
	if err := json.Unmarshal(body, &statisticsDTO); err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	if statisticsDTO.Max == nil || statisticsDTO.Min == nil || statisticsDTO.Average == nil || statisticsDTO.Sum == nil || statisticsDTO.HasDiagonalMatrix == nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	statistics := models.Statistics{
		Max:               *statisticsDTO.Max,
		Min:               *statisticsDTO.Min,
		Average:           *statisticsDTO.Average,
		Sum:               *statisticsDTO.Sum,
		HasDiagonalMatrix: *statisticsDTO.HasDiagonalMatrix,
	}

	return statistics, nil
}
