package config

import (
	"fmt"

	"github.com/isutare412/crawlert/internal/log"
	"github.com/isutare412/crawlert/internal/telegram"
)

type Config struct {
	Log    LogConfig    `koanf:"log"`
	Alerts AlertsConfig `koanf:"alerts"`
}

func (c Config) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("validating log config: %w", err)
	}
	if err := c.Alerts.Validate(); err != nil {
		return fmt.Errorf("validating alerts config: %w", err)
	}
	return nil
}

func (c Config) ToLogConfig() log.Config {
	return log.Config(c.Log)
}

func (c Config) ToTelegramMessageSenderConfigs() []telegram.MessageSenderConfig {
	cfgs := make([]telegram.MessageSenderConfig, 0, len(c.Alerts.Telegram.ChatIDs))
	for _, id := range c.Alerts.Telegram.ChatIDs {
		cfgs = append(cfgs, telegram.MessageSenderConfig{
			BotToken: c.Alerts.Telegram.BotToken,
			ChatID:   id,
		})
	}
	return cfgs
}

type LogConfig struct {
	Format log.Format `koanf:"format"`
	Level  log.Level  `koanf:"level"`
	Caller bool       `koanf:"caller"`
}

func (c LogConfig) Validate() error {
	if err := c.Format.Validate(); err != nil {
		return fmt.Errorf("validating format: %w", err)
	}
	if err := c.Level.Validate(); err != nil {
		return fmt.Errorf("validating level: %w", err)
	}
	return nil
}

type AlertsConfig struct {
	Telegram TelegramConfig `koanf:"telegram"`
}

func (c AlertsConfig) Validate() error {
	if err := c.Telegram.Validate(); err != nil {
		return fmt.Errorf("validating telegram config: %w", err)
	}
	return nil
}

type TelegramConfig struct {
	BotToken string   `koanf:"bot-token"`
	ChatIDs  []string `koanf:"chat-ids"`
}

func (c TelegramConfig) Validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("bot token should not be empty")
	}
	if len(c.ChatIDs) == 0 {
		return fmt.Errorf("chat ids should not be empty")
	}
	return nil
}
