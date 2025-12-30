package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// CommandHandler is the interface for command handlers.
type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// ApplyCommandDecorators applies the command decorators to the given handler.
func ApplyCommandDecorators[H, R any](handler CommandHandler[H, R], logger *logrus.Entry, client MetricsClient) CommandHandler[H, R] {
	return commandLoggingDecorator[H, R]{
		logger: logger,
		base: commandMetricsDecorator[H, R]{
			base:   handler,
			client: client,
		},
	}
}
