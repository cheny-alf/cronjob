package cronjob

import (
	"context"
	"time"
)

type AsyncLoadConfigParam struct {
	Interval int
	CallBack func(context.Context)
}

// AsyncLoadConfig synchronize the latest config
func AsyncLoadConfig(ctx context.Context, param *AsyncLoadConfigParam) {
	ticker := time.NewTicker(time.Minute * time.Duration(param.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go param.CallBack(ctx)
		case <-ctx.Done():
			return
		}
	}
}
