package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxMessageLength = 2000

type MessageSender struct {
	httpClient *http.Client

	webhookURL string
}

func NewMessageSender(cfg MessageSenderConfig) *MessageSender {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 100

	return &MessageSender{
		httpClient: &http.Client{Transport: transport},
		webhookURL: cfg.WebhookURL,
	}
}

func (s *MessageSender) SendMessage(ctx context.Context, message string) error {
	if len(message) == 0 {
		return nil
	}

	chunks := splitMessage(message)
	for i, chunk := range chunks {
		if err := s.sendChunk(ctx, chunk); err != nil {
			return fmt.Errorf("sending chunk %d/%d: %w", i+1, len(chunks), err)
		}
	}

	return nil
}

func (s *MessageSender) sendChunk(ctx context.Context, content string) error {
	reqBody, err := json.Marshal(webhookRequest{Content: content})
	if err != nil {
		return fmt.Errorf("marshaling message body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.webhookURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("building http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("doing http request: %w", err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %s; body (%s)", httpResp.Status, string(bodyBytes))
	}

	return nil
}
