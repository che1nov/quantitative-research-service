package config

import "os"

// Config хранит параметры запуска приложения.
type Config struct {
	Port             string
	LogLevel         string
	CSRFToken        string
	ExternalAPIToken string
	PublicBaseURL    string
	CORSAllowOrigin  string
	JWTSecret        string
}

// Load загружает конфигурацию из переменных окружения.
func Load() Config {
	port := getEnv("APP_PORT", "8080")

	return Config{
		Port:             port,
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		CSRFToken:        getEnv("CSRF_TOKEN", "dev-csrf-token"),
		ExternalAPIToken: getEnv("EXTERNAL_API_TOKEN", "dev-external-token"),
		PublicBaseURL:    getEnv("PUBLIC_BASE_URL", "http://localhost:"+port),
		CORSAllowOrigin:  getEnv("CORS_ALLOW_ORIGIN", "*"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-jwt-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
