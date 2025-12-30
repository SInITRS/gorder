package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
)

// MetricsClient is the interface for metrics clients.
type MetricsClient interface {
	Inc(key string, value int)
}

// queryMetricsDecorator is a decorator that measures the query execution time.
type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

// Handle implements QueryHandler.
func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	logger.Debug("Executing query")
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("querys.%s.durationh", actionName), int(end.Seconds()))
		if err != nil {
			q.client.Inc(fmt.Sprintf("querys.%s.errorh", actionName), int(end.Seconds()))
		} else {
			q.client.Inc(fmt.Sprintf("querys.%s.successh", actionName), int(end.Seconds()))
		}
	}()
	return q.base.Handle(ctx, cmd)
}

// commandMetricsDecorator is a decorator that measures the command execution time.
type commandMetricsDecorator[C, R any] struct {
	base   CommandHandler[C, R]
	client MetricsClient
}

// Handle implements CommandHandler.
func (c commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	logger.Debug("Executing command")
	defer func() {
		end := time.Since(start)
		c.client.Inc(fmt.Sprintf("commands.%s.durationh", actionName), int(end.Seconds()))
		if err != nil {
			c.client.Inc(fmt.Sprintf("commands.%s.errorh", actionName), int(end.Seconds()))
		} else {
			c.client.Inc(fmt.Sprintf("commands.%s.successh", actionName), int(end.Seconds()))
		}
	}()
	return c.base.Handle(ctx, cmd)
}
