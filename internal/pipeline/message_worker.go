package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/isutare412/crawlert/internal/core/model"
	"github.com/isutare412/crawlert/internal/core/port"
	"github.com/isutare412/crawlert/internal/log"
)

var regexPatternVariable = regexp.MustCompile(`\$\{?(\w+)\}?`)

type messageWorker struct {
	template       string
	messageSenders []port.MessageSender
	queryOutputs   <-chan queryOutput
	wg             sync.WaitGroup
}

func newMessageWorker(
	message string,
	messageSenders []port.MessageSender,
	queryOutputs <-chan queryOutput,
) *messageWorker {
	return &messageWorker{
		template:       message,
		messageSenders: messageSenders,
		queryOutputs:   queryOutputs,
		wg:             sync.WaitGroup{},
	}
}

func (w *messageWorker) run() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer log.RecoverIfPanic()

		for output := range w.queryOutputs {
			ctx := output.ctx

			if err := w.sendMessage(ctx, output.queryResult); err != nil {
				slog.ErrorContext(ctx, "failed to send message", "error", err)
			}
			slog.InfoContext(ctx, "sent message as query matched")
		}
	}()
}

func (w *messageWorker) shutdown() {
	w.wg.Wait()
}

func (w *messageWorker) sendMessage(
	ctx context.Context,
	queryRes model.QueryResult,
) error {
	message := buildMessage(w.template, queryRes.Variables)

	eg := errgroup.Group{}
	for _, sender := range w.messageSenders {
		eg.Go(func() error {
			if err := sender.SendMessage(ctx, message); err != nil {
				return fmt.Errorf("sending message: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func buildMessage(template string, variables map[string]string) string {
	return regexPatternVariable.ReplaceAllStringFunc(template, func(match string) string {
		key := strings.Trim(match, "${}")

		if value, ok := variables[key]; ok {
			return value
		}
		return match
	})
}
