package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiSendMessage = "/sendMessage"
)

type MessageSender struct {
	httpClient *http.Client

	botToken string
	chatID   string
}

func NewMessageSender(cfg Config) *MessageSender {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 100

	return &MessageSender{
		httpClient: &http.Client{Transport: transport},
		botToken:   cfg.BotToken,
		chatID:     cfg.ChatID,
	}
}

func (s *MessageSender) SendMessage(ctx context.Context, message string) error {
	if len(message) == 0 {
		return nil
	}
	message = truncateLargeMessage(message)

	req := sendMessageRequest{
		ChatID: s.chatID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return fmt.Errorf("marshaling message body: %w", err)
	}

	url := telegramBotAPIBase(s.botToken) + apiSendMessage
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return fmt.Errorf("building http request: %w", err)
	}

	httpReq.Header.Add("Content-Type", "application/json")

	httpResp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("doing http request: %w", err)
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status code %s; body (%s)", httpResp.Status, string(bodyBytes))
	}

	return nil
}

func truncateLargeMessage(msg string) string {
	// Telegram Bot API supports up to 4096 characters.
	// https://core.telegram.org/bots/api#sendmessage
	if len(msg) > 4096 {
		return msg[:4096]
	}
	return msg
}

func telegramBotAPIBase(botToken string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s", botToken)
}
