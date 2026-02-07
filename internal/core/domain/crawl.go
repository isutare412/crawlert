package domain

import (
	"net/http"
)

type CrawlRequest struct {
	URL    string
	Method string
	Header http.Header
	Body   []byte
}

type CrawlResponse struct {
	Header http.Header
	Body   []byte
}
