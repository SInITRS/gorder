package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func ApplyCommandDecorators[H, R any](handler CommandHandler[H, R], logger *logrus.Entry, client MetricsClient) CommandHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: client,
		},
	}
}
