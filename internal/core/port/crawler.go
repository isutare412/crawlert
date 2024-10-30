package port

import (
	"context"

	"github.com/isutare412/crawlert/internal/core/model"
)

type HTTPCrawler interface {
	Crawl(context.Context, model.CrawlRequest) (model.CrawlResponse, error)
}
