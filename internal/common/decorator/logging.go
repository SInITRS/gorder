package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// queryLoggingDecorator is a decorator that logs the query execution.
type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

// Handle implements QueryHandler.
func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing query")
	defer func() {
		if err != nil {
			logger.Error("Failed to execute query", err)
		} else {
			logger.Debug("Successfully executed query")
		}
	}()
	return q.base.Handle(ctx, cmd)
}

// commandLoggingDecorator is a decorator that logs the command execution.
type commandLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   CommandHandler[C, R]
}

// Handle implements CommandHandler.
func (c commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := c.logger.WithFields(logrus.Fields{
		"command":      generateActionName(cmd),
		"command_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing command")
	defer func() {
		if err != nil {
			logger.Error("Failed to execute command", err)
		} else {
			logger.Debug("Successfully executed command")
		}
	}()
	return c.base.Handle(ctx, cmd)
}

func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
