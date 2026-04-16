package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	NodeName          string
	Runtime           string
	Version           string
	Capacity          int
	APIServerEndpoint string
	HeartbeatInterval time.Duration
	HealthListen      string
	ShutdownTimeout   time.Duration
}

func Load() Config {
	return Config{
		NodeName:          getEnv("MUBELET_NODE_NAME", hostnameOrFallback()),
		Runtime:           getEnv("MUBELET_RUNTIME", "containerd"),
		Version:           getEnv("MUBELET_VERSION", "0.1.0"),
		Capacity:          getEnvInt("MUBELET_CAPACITY", 10),
		APIServerEndpoint: getEnv("MUBE_API_SERVER", "http://127.0.0.1:8080"),
		HeartbeatInterval: getEnvDuration("MUBELET_HEARTBEAT_INTERVAL", 10*time.Second),
		HealthListen:      getEnv("MUBELET_HEALTH_LISTEN", ":10250"),
		ShutdownTimeout:   getEnvDuration("MUBELET_SHUTDOWN_TIMEOUT", 5*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	if n < 1 {
		return fallback
	}
	return n
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	if d <= 0 {
		return fallback
	}
	return d
}

func hostnameOrFallback() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "mube-node-1"
	}
	return h
}
