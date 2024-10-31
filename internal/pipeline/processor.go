package pipeline

import (
	"fmt"
	"log/slog"

	"github.com/isutare412/crawlert/internal/core/port"
)

type Processor struct {
	workerGroups []*workerGroup
}

func NewProcessor(
	cfg ProcessorConfig,
	httpCrawler port.HTTPCrawler,
	messageSenders []port.MessageSender,
) (*Processor, error) {
	cfgsEnabled := filterEnabledConfig(cfg.Crawls)
	if len(cfgsEnabled) == 0 {
		return nil, fmt.Errorf("all crawls are disabled")
	}

	workerGroups := make([]*workerGroup, 0, len(cfgsEnabled))
	for _, cfg := range cfgsEnabled {
		group, err := newWorkerGroup(cfg, httpCrawler, messageSenders)
		if err != nil {
			return nil, fmt.Errorf("creating worker group of %s: %w", cfg.Name, err)
		}

		slog.Info("worker group created", "jobName", cfg.Name, "interval", cfg.Interval.String())
		workerGroups = append(workerGroups, group)
	}

	return &Processor{
		workerGroups: workerGroups,
	}, nil
}

func (p *Processor) Run() {
	for _, group := range p.workerGroups {
		group.run()
	}
}

func (p *Processor) Shutdown() {
	for _, group := range p.workerGroups {
		group.shutdown()
	}
}

func filterEnabledConfig(cfgs []CrawlConfig) []CrawlConfig {
	enabled := make([]CrawlConfig, 0, len(cfgs))
	for _, cfg := range cfgs {
		if cfg.Enabled {
			enabled = append(enabled, cfg)
		}
	}
	return enabled
}
