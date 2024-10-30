package config

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/isutare412/crawlert/internal/log"
	"github.com/isutare412/crawlert/internal/telegram"
)

type Config struct {
	Log    LogConfig     `koanf:"log"`
	Crawls []CrawlConfig `koanf:"crawls"`
	Alerts AlertsConfig  `koanf:"alerts"`
}

func (c Config) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("validating log config: %w", err)
	}
	if err := c.Alerts.Validate(); err != nil {
		return fmt.Errorf("validating alerts config: %w", err)
	}

	for _, cfg := range c.Crawls {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validating crawl config: %w", err)
		}
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

type CrawlConfig struct {
	Name     string            `koanf:"name"`
	Enabled  bool              `koanf:"enabled"`
	Interval time.Duration     `koanf:"interval"`
	Target   CrawlTargetConfig `koanf:"target"`
	Query    CrawlQueryConfig  `koanf:"query"`
	Message  string            `koanf:"message"`
}

func (c CrawlConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name should not be empty")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("interval %v of %s should not be empty or negative", c.Name, c.Interval)
	}
	if c.Message == "" {
		return fmt.Errorf("message of %s should not be empty", c.Name)
	}

	if err := c.Target.Validate(); err != nil {
		return fmt.Errorf("validating target config of %s: %w", c.Name, err)
	}

	return nil
}

type CrawlTargetConfig struct {
	HTTP CrawlHTTPTargetConfig `koanf:"http"`
}

func (c CrawlTargetConfig) Validate() error {
	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("validating http target: %w", err)
	}

	return nil
}

type CrawlHTTPTargetConfig struct {
	Method string            `koanf:"method"`
	URL    string            `koanf:"url"`
	Header map[string]string `koanf:"header"`
	Body   string            `koanf:"body"`
}

func (c CrawlHTTPTargetConfig) Validate() error {
	switch c.Method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
	default:
		return fmt.Errorf("unexpected mdethod %s", c.Method)
	}

	if _, err := url.Parse(c.URL); err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	return nil
}

type CrawlQueryConfig struct {
	Check     string            `koanf:"check"`
	Variables map[string]string `koanf:"variables"`
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
