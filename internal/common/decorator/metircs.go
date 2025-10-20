package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/gopkg/util/logger"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

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
