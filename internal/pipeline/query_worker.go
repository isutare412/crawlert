package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/isutare412/crawlert/internal/core/model"
	"github.com/isutare412/crawlert/internal/core/port"
	"github.com/isutare412/crawlert/internal/log"
	"github.com/isutare412/crawlert/internal/query"
)

type queryWorker struct {
	applier      port.QueryApplier
	crawlOutputs <-chan crawlOutput
	queryOutputs chan<- queryOutput
	wg           sync.WaitGroup
}

func newQueryWorker(
	cfg CrawlQueryConfig,
	crawlOutputs <-chan crawlOutput,
	queryOutputs chan<- queryOutput,
) (*queryWorker, error) {
	applier, err := query.NewApplier(cfg.Check, cfg.Variables)
	if err != nil {
		return nil, fmt.Errorf("creating query applier: %w", err)
	}

	return &queryWorker{
		applier:      applier,
		crawlOutputs: crawlOutputs,
		queryOutputs: queryOutputs,
		wg:           sync.WaitGroup{},
	}, nil
}

func (w *queryWorker) run() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer log.RecoverIfPanic()

		for output := range w.crawlOutputs {
			ctx := output.ctx

			queryResult, err := w.query(ctx, output.crawlResponse)
			switch {
			case err != nil:
				slog.ErrorContext(ctx, "failed to apply query", "error", err)
				continue
			case !queryResult.Matched:
				slog.InfoContext(ctx, "query result does not matched")
				continue
			}

			w.queryOutputs <- queryOutput{
				ctx:           ctx,
				crawlResponse: output.crawlResponse,
				queryResult:   queryResult,
			}
		}
	}()
}

func (w *queryWorker) shutdown() {
	w.wg.Wait()
}

func (w *queryWorker) query(ctx context.Context, crawlResp model.CrawlResponse) (model.QueryResult, error) {
	result, err := w.applier.ApplyQuery(crawlResp.Body)
	if err != nil {
		return model.QueryResult{}, fmt.Errorf("applying query: %w", err)
	}

	return result, nil
}
