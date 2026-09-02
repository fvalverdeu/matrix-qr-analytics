package configs

import "testing"

func setLoadBaselineEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PORT", "")
	t.Setenv("STATISTICS_SERVICE_URL", "http://localhost:3000")
	t.Setenv("STATISTICS_TIMEOUT_MS", "")
	t.Setenv("STATISTICS_AUTH_MODE", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
}

func TestLoad_DefaultStatisticsAuthModeIsNone(t *testing.T) {
	setLoadBaselineEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.StatisticsAuthMode != StatisticsAuthModeNone {
		t.Fatalf("expected default auth mode %q, got %q", StatisticsAuthModeNone, cfg.StatisticsAuthMode)
	}
}

func TestLoad_ExplicitNoneAuthMode(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("STATISTICS_AUTH_MODE", "none")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.StatisticsAuthMode != StatisticsAuthModeNone {
		t.Fatalf("expected auth mode %q, got %q", StatisticsAuthModeNone, cfg.StatisticsAuthMode)
	}
}

func TestLoad_ExplicitGoogleIDTokenAuthMode(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("STATISTICS_AUTH_MODE", "google-id-token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.StatisticsAuthMode != StatisticsAuthModeGoogleIDToken {
		t.Fatalf("expected auth mode %q, got %q", StatisticsAuthModeGoogleIDToken, cfg.StatisticsAuthMode)
	}
}

func TestLoad_UnsupportedStatisticsAuthModeFailsFast(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("STATISTICS_AUTH_MODE", "oauth")

	_, err := Load()
	if err == nil {
		t.Fatal("expected configuration error for unsupported STATISTICS_AUTH_MODE")
	}
}

func TestLoad_DefaultCORSAllowedOrigins(t *testing.T) {
	setLoadBaselineEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.CORSAllowedOrigins != DefaultCORSAllowedOrigins {
		t.Fatalf("expected default CORS origins %q, got %q", DefaultCORSAllowedOrigins, cfg.CORSAllowedOrigins)
	}
}

func TestLoad_ExplicitCORSAllowedOriginsAreTrimmed(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", " https://app.example.com, http://localhost:5173 ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	const expected = "https://app.example.com,http://localhost:5173"
	if cfg.CORSAllowedOrigins != expected {
		t.Fatalf("expected CORS origins %q, got %q", expected, cfg.CORSAllowedOrigins)
	}
}

func TestLoad_CORSAllowedOriginsRejectsWildcard(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "*")

	_, err := Load()
	if err == nil {
		t.Fatal("expected configuration error for wildcard CORS origin")
	}
}

func TestLoad_CORSAllowedOriginsRejectsEmptyEntries(t *testing.T) {
	setLoadBaselineEnv(t)
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173, ")

	_, err := Load()
	if err == nil {
		t.Fatal("expected configuration error for empty CORS origin entry")
	}
}
