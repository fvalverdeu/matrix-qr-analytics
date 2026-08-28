package configs

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	StatisticsServiceURL string
	StatisticsTimeout    time.Duration
}

func Load() Config {
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

	return Config{
		Port:                 port,
		StatisticsServiceURL: os.Getenv("STATISTICS_SERVICE_URL"),
		StatisticsTimeout:    time.Duration(timeoutMS) * time.Millisecond,
	}
}
