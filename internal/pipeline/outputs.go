package pipeline

import (
	"context"

	"github.com/isutare412/crawlert/internal/core/model"
)

type triggerOutput struct {
	ctx          context.Context
	crawlRequest model.CrawlRequest
}

type crawlOutput struct {
	ctx           context.Context
	crawlResponse model.CrawlResponse
}

type queryOutput struct {
	ctx           context.Context
	crawlResponse model.CrawlResponse
	queryResult   model.QueryResult
}

type messageOutput struct {
	ctx     context.Context
	message string
}
