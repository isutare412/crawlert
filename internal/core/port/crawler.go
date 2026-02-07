package port

import (
	"context"

	"github.com/isutare412/crawlert/internal/core/domain"
)

type HTTPCrawler interface {
	Crawl(context.Context, domain.CrawlRequest) (domain.CrawlResponse, error)
}
