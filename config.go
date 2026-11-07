package sieve

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileConfig is the YAML shape loaded from disk for the sidecar and library users.
type FileConfig struct {
	BlockOn        string       `yaml:"block_on"`
	IncludeInput   bool         `yaml:"include_input"`
	DisableBuiltin bool         `yaml:"disable_builtin"`
	LogPath        string       `yaml:"log_path"`
	WebhookURL     string       `yaml:"webhook_url"`
	ListenAddr     string       `yaml:"listen_addr"`
	Rules          []FileRule   `yaml:"rules"`
}

// FileRule is a custom detection rule from YAML.
type FileRule struct {
	Category string `yaml:"category"`
	Pattern  string `yaml:"pattern"`
	Risk     string `yaml:"risk"`
}

// LoadFile reads and parses a YAML config file.
func LoadFile(path string) (FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileConfig{}, err
	}
	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return FileConfig{}, fmt.Errorf("parse config: %w", err)
	}
	cfg.ListenAddr = resolveListenAddr(cfg.ListenAddr)
	if cfg.BlockOn == "" {
		cfg.BlockOn = string(RiskBlocked)
	}
	if cfg.LogPath == "" {
		cfg.LogPath = "sieve.log"
	}
	return cfg, nil
}

// resolveListenAddr prefers Railway/Heroku-style $PORT when set, otherwise
// the YAML listen_addr, otherwise :8080. Always binds on all interfaces.
func resolveListenAddr(configured string) string {
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		if strings.HasPrefix(port, ":") {
			return "0.0.0.0" + port
		}
		return "0.0.0.0:" + port
	}
	if configured == "" {
		return ":8080"
	}
	return configured
}

// ToConfig converts file config into a runtime Config.
func (f FileConfig) ToConfig() Config {
	custom := make([]Rule, 0, len(f.Rules))
	for _, r := range f.Rules {
		custom = append(custom, Rule{
			Category: r.Category,
			Pattern:  r.Pattern,
			Risk:     Risk(r.Risk),
		})
	}
	return Config{
		BlockOn:        Risk(f.BlockOn),
		IncludeInput:   f.IncludeInput,
		DisableBuiltin: f.DisableBuiltin,
		CustomRules:    custom,
	}
}
