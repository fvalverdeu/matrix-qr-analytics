package clients

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"qr-service-go/internal/models"
)

type HTTPStatisticsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPStatisticsClient(baseURL string, timeout time.Duration) *HTTPStatisticsClient {
	return &HTTPStatisticsClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
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

	var statistics models.Statistics
	if err := json.Unmarshal(body, &statistics); err != nil {
		return models.Statistics{}, newUnavailableError("Statistics service is unavailable")
	}

	return statistics, nil
}
