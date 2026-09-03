package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	StatisticsAuthModeNone          = "none"
	StatisticsAuthModeGoogleIDToken = "google-id-token"
	DefaultCORSAllowedOrigins       = "http://localhost:5173"
)

type Config struct {
	Port                 string
	StatisticsServiceURL string
	StatisticsTimeout    time.Duration
	StatisticsAuthMode   string
	CORSAllowedOrigins   string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	timeoutMS := 5000
	if value := os.Getenv("STATISTICS_TIMEOUT_MS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			timeoutMS = parsed
		}
	}

	authMode, err := parseStatisticsAuthMode(os.Getenv("STATISTICS_AUTH_MODE"))
	if err != nil {
		return Config{}, err
	}

	allowedOrigins, err := parseCORSAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:                 port,
		StatisticsServiceURL: os.Getenv("STATISTICS_SERVICE_URL"),
		StatisticsTimeout:    time.Duration(timeoutMS) * time.Millisecond,
		StatisticsAuthMode:   authMode,
		CORSAllowedOrigins:   allowedOrigins,
	}, nil
}

func parseStatisticsAuthMode(value string) (string, error) {
	normalized := strings.TrimSpace(strings.ToLower(value))
	if normalized == "" {
		return StatisticsAuthModeNone, nil
	}

	switch normalized {
	case StatisticsAuthModeNone, StatisticsAuthModeGoogleIDToken:
		return normalized, nil
	default:
		return "", fmt.Errorf("invalid STATISTICS_AUTH_MODE %q: supported values are %q and %q", value, StatisticsAuthModeNone, StatisticsAuthModeGoogleIDToken)
	}
}

func parseCORSAllowedOrigins(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return DefaultCORSAllowedOrigins, nil
	}

	rawOrigins := strings.Split(normalized, ",")
	origins := make([]string, 0, len(rawOrigins))
	for _, rawOrigin := range rawOrigins {
		origin := strings.TrimSpace(rawOrigin)
		if origin == "" {
			return "", fmt.Errorf("invalid CORS_ALLOWED_ORIGINS %q: empty origin entries are not allowed", value)
		}
		if origin == "*" {
			return "", fmt.Errorf("invalid CORS_ALLOWED_ORIGINS %q: wildcard origin %q is not allowed", value, origin)
		}
		origins = append(origins, origin)
	}

	return strings.Join(origins, ","), nil
}
