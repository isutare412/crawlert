package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/isutare412/crawlert/internal/core/model"
)

type Crawler struct {
	client *http.Client
}

func NewCrawler() *Crawler {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 100

	return &Crawler{
		client: &http.Client{Transport: transport},
	}
}

func (c *Crawler) Crawl(ctx context.Context, req model.CrawlRequest) (model.CrawlResponse, error) {
	bodyBuffer := bytes.NewBuffer(req.Body)
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyBuffer)
	if err != nil {
		return model.CrawlResponse{}, fmt.Errorf("creating http request: %w", err)
	}

	for key, values := range req.Header {
		for _, v := range values {
			httpReq.Header.Add(key, v)
		}
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return model.CrawlResponse{}, fmt.Errorf("doing http request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode >= http.StatusBadRequest {
		return model.CrawlResponse{}, fmt.Errorf("unexpected http response code '%s'", httpResp.Status)
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return model.CrawlResponse{}, fmt.Errorf("reading http response body: %w", err)
	}

	return model.CrawlResponse{
		Header: httpResp.Header.Clone(),
		Body:   bodyBytes,
	}, nil
}
