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

	var messageSenders []*telegram.MessageSender
	for _, cfg := range cfg.ToTelegramMessageSenderConfigs() {
		sender := telegram.NewMessageSender(cfg)
		messageSenders = append(messageSenders, sender)
	}

	pipelineProcessor, err := pipeline.NewProcessor(
		cfg.ToPipelineProcessorConfig(),
		httpCrawler,
		toMessageSenderInterfaces(messageSenders))
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

func toMessageSenderInterfaces(s []*telegram.MessageSender) []port.MessageSender {
	res := make([]port.MessageSender, 0, len(s))
	for _, ms := range s {
		res = append(res, ms)
	}
	return res
}
