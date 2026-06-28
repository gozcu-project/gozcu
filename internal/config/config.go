package config

import (
	"os"
	"time"
)

type Config struct {
	BackendBaseURL string
	CACertPath     string
	ClientCertPath string
	ClientKeyPath  string
	WaitTimeout    time.Duration
}

func Load() Config {
	return Config{
		BackendBaseURL: getEnv("GOZCU_BACKEND_URL", "https://localhost:8443/api/approvals"),
		CACertPath:     getEnv("GOZCU_CA_CERT", "/opt/gozcu/etc/pki/ca.crt"),
		ClientCertPath: getEnv("GOZCU_CLIENT_CERT", "/opt/gozcu/etc/pki/gate.crt"),
		ClientKeyPath:  getEnv("GOZCU_CLIENT_KEY", "/opt/gozcu/etc/pki/gate.key"),
		WaitTimeout:    30 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}