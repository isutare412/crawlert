package port

import "context"

type MessageSender interface {
	SendMessage(ctx context.Context, message string) error
}
