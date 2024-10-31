package pipeline

import (
	"time"
)

type ProcessorConfig struct {
	Crawls []CrawlConfig
}

type CrawlConfig struct {
	Name     string
	Enabled  bool
	Interval time.Duration
	Target   CrawlTargetConfig
	Query    CrawlQueryConfig
	Message  string
}

type CrawlTargetConfig struct {
	HTTP CrawlHTTPTargetConfig
}

type CrawlHTTPTargetConfig struct {
	Method string
	URL    string
	Header map[string]string
	Body   string
}

type CrawlQueryConfig struct {
	Check     string
	Variables map[string]string
}
