package pipeline

import (
	"context"

	"github.com/isutare412/crawlert/internal/core/domain"
)

type triggerOutput struct {
	ctx          context.Context
	crawlRequest domain.CrawlRequest
}

type crawlOutput struct {
	ctx           context.Context
	crawlResponse domain.CrawlResponse
}

type queryOutput struct {
	ctx           context.Context
	crawlResponse domain.CrawlResponse
	queryResult   domain.QueryResult
}

type messageOutput struct {
	ctx     context.Context
	message string
}
