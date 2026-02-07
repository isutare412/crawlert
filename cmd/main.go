package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/isutare412/crawlert/internal/config"
	"github.com/isutare412/crawlert/internal/core/port"
	"github.com/isutare412/crawlert/internal/discord"
	"github.com/isutare412/crawlert/internal/http"
	"github.com/isutare412/crawlert/internal/log"
	"github.com/isutare412/crawlert/internal/pipeline"
	"github.com/isutare412/crawlert/internal/telegram"
)

var configPath = flag.String("configs", ".", "path to config directory")

func init() {
	flag.Parse()
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	log.Init(cfg.ToLogConfig())

	httpCrawler := http.NewCrawler()

	messageSenders, err := buildMessageSenders(cfg)
	if err != nil {
		slog.Error("failed to build message senders", "error", err)
		os.Exit(1)
	}

	pipelineProcessor, err := pipeline.NewProcessor(
		cfg.ToPipelineProcessorConfig(),
		httpCrawler,
		messageSenders)
	if err != nil {
		slog.Error("failed to create pipeline processor", "error", err)
		os.Exit(1)
	}

	pipelineProcessor.Run()
	waitUntilSignal()
	pipelineProcessor.Shutdown()
}

func waitUntilSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	s := <-signals
	slog.Info("received signal", "signal", s.String())
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(*configPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	slog.Debug("loaded config", "config", cfg)
	return cfg, nil
}

func buildMessageSenders(cfg *config.Config) ([]port.MessageSender, error) {
	switch cfg.Alerts.Type {
	case "telegram":
		var senders []port.MessageSender
		for _, c := range cfg.ToTelegramMessageSenderConfigs() {
			senders = append(senders, telegram.NewMessageSender(c))
		}
		return senders, nil
	case "discord":
		var senders []port.MessageSender
		for _, c := range cfg.ToDiscordMessageSenderConfigs() {
			senders = append(senders, discord.NewMessageSender(c))
		}
		return senders, nil
	default:
		return nil, fmt.Errorf("unknown alerts type: %s", cfg.Alerts.Type)
	}
}
