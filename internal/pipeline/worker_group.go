package pipeline

import (
	"fmt"

	"github.com/isutare412/crawlert/internal/core/port"
)

type workerGroup struct {
	trigger *triggerWorker
	crawl   *crawlWorker
	query   *queryWorker
	message *messageWorker

	triggerOutputs chan triggerOutput
	crawlOutputs   chan crawlOutput
	queryOutputs   chan queryOutput
}

func newWorkerGroup(
	cfg CrawlConfig,
	httpCrawler port.HTTPCrawler,
	messageSenders []port.MessageSender,
) (*workerGroup, error) {
	var (
		triggerOutputs = make(chan triggerOutput, 1)
		crawlOutputs   = make(chan crawlOutput, 1)
		queryOutputs   = make(chan queryOutput, 1)
	)

	triggerWorker := newTriggerWorker(cfg, triggerOutputs)
	crawlWorker := newCrawlWorker(httpCrawler, triggerOutputs, crawlOutputs)

	queryWorker, err := newQueryWorker(cfg.Query, crawlOutputs, queryOutputs)
	if err != nil {
		return nil, fmt.Errorf("creating query worker: %w", err)
	}

	messageWorker := newMessageWorker(cfg.Message, messageSenders, queryOutputs)

	return &workerGroup{
		trigger:        triggerWorker,
		crawl:          crawlWorker,
		query:          queryWorker,
		message:        messageWorker,
		triggerOutputs: triggerOutputs,
		crawlOutputs:   crawlOutputs,
		queryOutputs:   queryOutputs,
	}, nil
}

func (g *workerGroup) run() {
	g.message.run()
	g.query.run()
	g.crawl.run()
	g.trigger.run()
}

func (g *workerGroup) shutdown() {
	g.trigger.shutdown()
	close(g.triggerOutputs)
	g.crawl.shutdown()
	close(g.crawlOutputs)
	g.query.shutdown()
	close(g.queryOutputs)
	g.message.shutdown()
}
