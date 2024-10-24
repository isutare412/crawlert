package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	slogmulti "github.com/samber/slog-multi"
)

func init() {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.RFC3339Nano,
		NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
	})

	logger := slog.New(
		slogmulti.
			Pipe(slogmulti.NewHandleInlineMiddleware(injectAttrsFromCtx)).
			Handler(handler),
	)
	slog.SetDefault(logger)
}

func Init(cfg Config) {
	var (
		writer    = os.Stdout
		level     = cfg.Level.SlogLevel()
		addSource = cfg.Caller
	)

	var handler slog.Handler
	switch cfg.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     level,
			AddSource: addSource,
		})
	case FormatText:
		handler = tint.NewHandler(writer, &tint.Options{
			Level:      level,
			TimeFormat: time.RFC3339,
			NoColor:    !isatty.IsTerminal(writer.Fd()),
			AddSource:  addSource,
		})
	}

	logger := slog.New(
		slogmulti.
			Pipe(slogmulti.NewHandleInlineMiddleware(injectAttrsFromCtx)).
			Handler(handler),
	)
	slog.SetDefault(logger)
}
