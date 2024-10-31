package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/isutare412/crawlert/internal/core/model"
	"github.com/isutare412/crawlert/internal/core/port"
	"github.com/isutare412/crawlert/internal/log"
)

type crawlWorker struct {
	httpCrawler    port.HTTPCrawler
	triggerOutputs <-chan triggerOutput
	crawlOutputs   chan<- crawlOutput
	wg             sync.WaitGroup
}

func newCrawlWorker(httpCrawler port.HTTPCrawler, triggerOutputs <-chan triggerOutput, crawlOutputs chan<- crawlOutput) *crawlWorker {
	return &crawlWorker{
		httpCrawler:    httpCrawler,
		triggerOutputs: triggerOutputs,
		crawlOutputs:   crawlOutputs,
		wg:             sync.WaitGroup{},
	}
}

func (w *crawlWorker) run() {
	w.wg.Add(1)
	go func() {
		w.wg.Done()
		defer log.RecoverIfPanic()

		for output := range w.triggerOutputs {
			ctx := output.ctx

			resp, err := w.crawl(ctx, output.crawlRequest)
			if err != nil {
				slog.ErrorContext(ctx, "failed to crawl", "error", err)
				continue
			}

			w.crawlOutputs <- crawlOutput{
				ctx:           ctx,
				crawlResponse: resp,
			}
		}
	}()
}

func (w *crawlWorker) shutdown() {
	w.wg.Wait()
}

func (w *crawlWorker) crawl(ctx context.Context, req model.CrawlRequest) (model.CrawlResponse, error) {
	resp, err := w.httpCrawler.Crawl(ctx, req)
	if err != nil {
		return model.CrawlResponse{}, fmt.Errorf("crawling http: %w", err)
	}

	return resp, nil
}
