package tgbotserver

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
)

// Mode represents available telegram bot server modes
type Mode string

const (
	none    Mode = "none"
	polling Mode = "polling"
	webhook Mode = "webhook"
)

type Config struct {
	Mode Mode `mapstructure:"mode"`
}

func (cfg *Config) InitDefault() {
	if string(cfg.Mode) == "" {
		cfg.Mode = none
	}
}

type PollingConfig struct {
	Offset         int      `mapstructure:"offset"`
	Limit          int      `mapstructure:"limit"`
	Timeout        int      `mapstructure:"timeout"`
	AllowedUpdates []string `mapstructure:"allowed_updates"`
}

func (cfg *PollingConfig) InitDefault() {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30
	}
}

func (cfg *PollingConfig) BuildUpdateConfig() tgbotapi.UpdateConfig {
	return tgbotapi.UpdateConfig{
		Offset:         cfg.Offset,
		Limit:          cfg.Limit,
		Timeout:        cfg.Timeout,
		AllowedUpdates: cfg.AllowedUpdates,
	}
}

type WebhookConfig struct {
	URL                string   `mapstructure:"url"`
	Certificate        string   `mapstructure:"certificate"`
	IPAddress          string   `mapstructure:"ip"`
	AllowedUpdates     []string `mapstructure:"allowed_updates"`
	MaxConnections     int      `mapstructure:"max_connections"`
	DropPendingUpdates bool     `mapstructure:"drop_pending_updates"`
}

func (cfg *WebhookConfig) BuildWebhookConfig() (tgbotapi.WebhookConfig, error) {
	u, err := url.Parse(cfg.URL)

	if err != nil {
		return tgbotapi.WebhookConfig{}, err
	}

	var cert tgbotapi.RequestFileData
	if cfg.Certificate != "" {
		cert = tgbotapi.FilePath(cfg.Certificate)
	}

	return tgbotapi.WebhookConfig{
		URL:                u,
		Certificate:        cert,
		IPAddress:          cfg.IPAddress,
		MaxConnections:     cfg.MaxConnections,
		AllowedUpdates:     cfg.AllowedUpdates,
		DropPendingUpdates: cfg.DropPendingUpdates,
	}, nil
}
