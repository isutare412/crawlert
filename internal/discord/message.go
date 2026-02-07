package discord

type webhookRequest struct {
	Content string `json:"content"`
}

// splitMessage splits a message into chunks that fit within Discord's
// 2000 character limit, preferring to split at newline boundaries.
func splitMessage(msg string) []string {
	if len(msg) <= maxMessageLength {
		return []string{msg}
	}

	var chunks []string
	for len(msg) > 0 {
		if len(msg) <= maxMessageLength {
			chunks = append(chunks, msg)
			break
		}

		chunk := msg[:maxMessageLength]
		splitAt := maxMessageLength

		if idx := lastIndexNewline(chunk); idx > 0 {
			splitAt = idx + 1
		}

		chunks = append(chunks, msg[:splitAt])
		msg = msg[splitAt:]
	}

	return chunks
}

func lastIndexNewline(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			return i
		}
	}
	return -1
}
