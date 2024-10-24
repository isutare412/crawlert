package log

import (
	"context"
	"log/slog"
)

// WithValue returns new context with key, val attached to ctx. The attached
// values are logged automatically by [slog.InfoContext], [slog.ErrorContext],
// ...
//
//	WithValue(context.Background(),
//		"foo", 123,
//		"bar", variable,
//	)
func WithValue(ctx context.Context, key string, val any) context.Context {
	bag, ok := getContextBag(ctx)
	if ok {
		bag.values[key] = val
		return ctx
	}

	bag = &contextBag{
		values: map[string]any{
			key: val,
		},
	}
	return context.WithValue(ctx, contextBagKey{}, bag)
}

// WithoutValue returns new context with keys removed. If the key does not
// exist, it does nothing.
func WithoutValue(ctx context.Context, keys ...string) context.Context {
	bag, ok := getContextBag(ctx)
	if !ok {
		return ctx
	}

	for _, k := range keys {
		delete(bag.values, k)
	}
	return ctx
}

type contextBagKey struct{}

type contextBag struct {
	values map[string]any
}

func getContextBag(ctx context.Context) (*contextBag, bool) {
	bag, ok := ctx.Value(contextBagKey{}).(*contextBag)
	return bag, ok
}

// injectAttrsFromCtx is a middleware function that read values from context and
// add them as [slog.Attr].
func injectAttrsFromCtx(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
	bag, ok := getContextBag(ctx)
	if ok {
		for k, v := range bag.values {
			record.AddAttrs(slog.Any(k, v))
		}
	}

	return next(ctx, record)
}
