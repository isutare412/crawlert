package telegram

import (
	"regexp"
)

var regexNeedEscape = regexp.MustCompile(`[\*\[\]\(\)\~\>\#\+\-\=\|\{\}\.\!]`)

type sendMessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func newSendMessageRequest(chatID, text string) sendMessageRequest {
	text = truncateLargeMessage(text)
	text = escapeForMarkdown(text)

	return sendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}
}

// truncateLargeMessage truncates msg if it is too long. Telegram Bot API
// supports up to 4096 characters.
// https://core.telegram.org/bots/api#sendmessage
func truncateLargeMessage(msg string) string {
	if len(msg) > 4096 {
		return msg[:4096]
	}
	return msg
}

// escapeForMarkdown escapes any character for Telegram MarkdownV2 parse mode.
// https://core.telegram.org/bots/api#markdownv2-style
func escapeForMarkdown(s string) string {
	return regexNeedEscape.ReplaceAllStringFunc(s, func(match string) string {
		return `\` + match
	})
}
