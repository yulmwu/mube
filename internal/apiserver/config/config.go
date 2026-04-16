package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddress       string `yaml:"listenAddress"`
	ShutdownTimeout     time.Duration
	NodeNotReadyTimeout time.Duration
	Nodes               []NodeConfig `yaml:"nodes"`
}

type NodeConfig struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func Load() (Config, error) {
	path := getEnv("MUBE_APISERVER_CONFIG", "configs/apiserver.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read apiserver config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse apiserver config %q: %w", path, err)
	}

	if cfg.ListenAddress == "" {
		cfg.ListenAddress = ":8080"
	}
	cfg.ShutdownTimeout = getEnvDuration("MUBE_APISERVER_SHUTDOWN_TIMEOUT", 5*time.Second)
	cfg.NodeNotReadyTimeout = getEnvDuration("MUBE_NODE_NOTREADY_TIMEOUT", 30*time.Second)

	if len(cfg.Nodes) == 0 {
		return Config{}, fmt.Errorf("apiserver config %q has no nodes", path)
	}

	for i := range cfg.Nodes {
		n := cfg.Nodes[i]
		if n.Name == "" {
			return Config{}, fmt.Errorf("apiserver config %q has node with empty name", path)
		}

		if n.IP == "" {
			return Config{}, fmt.Errorf("apiserver config %q has node %q with empty ip", path, n.Name)
		}

		if n.Port <= 0 || n.Port > 65535 {
			return Config{}, fmt.Errorf("apiserver config %q has node %q with invalid port %d", path, n.Name, n.Port)
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
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
