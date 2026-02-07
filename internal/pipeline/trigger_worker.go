package pipeline

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/isutare412/crawlert/internal/core/domain"
	"github.com/isutare412/crawlert/internal/log"
)

type triggerWorker struct {
	jobName        string
	interval       time.Duration
	crawlRequest   domain.CrawlRequest
	triggerOutputs chan<- triggerOutput

	lifetimeCtx    context.Context
	lifetimeCancel context.CancelFunc
	wg             sync.WaitGroup
}

func newTriggerWorker(cfg CrawlConfig, triggerOutputs chan<- triggerOutput) *triggerWorker {
	ctx, cancel := context.WithCancel(context.Background())

	return &triggerWorker{
		jobName:        cfg.Name,
		interval:       cfg.Interval,
		crawlRequest:   buildCrawlRequest(cfg.Target.HTTP),
		triggerOutputs: triggerOutputs,
		lifetimeCtx:    ctx,
		lifetimeCancel: cancel,
		wg:             sync.WaitGroup{},
	}
}

func (w *triggerWorker) run() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer log.RecoverIfPanic()

		for {
			w.trigger()

			select {
			case <-time.After(w.interval):
			case <-w.lifetimeCtx.Done():
				return
			}
		}
	}()
}

func (w *triggerWorker) shutdown() {
	w.lifetimeCancel()
	w.wg.Wait()
}

func (w *triggerWorker) trigger() {
	ctx := context.Background()
	ctx = log.WithValue(ctx, "jobName", w.jobName)

	w.triggerOutputs <- triggerOutput{
		ctx:          ctx,
		crawlRequest: w.crawlRequest,
	}
}

func buildCrawlRequest(cfg CrawlHTTPTargetConfig) domain.CrawlRequest {
	return domain.CrawlRequest{
		URL:    cfg.URL,
		Method: cfg.Method,
		Header: buildHTTPHeader(cfg.Header),
		Body:   []byte(cfg.Body),
	}
}

func buildHTTPHeader(h map[string]string) http.Header {
	out := http.Header{}
	for k, v := range h {
		out.Set(k, v)
	}
	return out
}
