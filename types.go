package cronjob

import (
	"context"
	"runtime"
)

func GoRecover(ctx context.Context, log string, trace bool) {
	if err := recover(); err != nil {
		if trace {
			trace := make([]byte, 4096)
			runtime.Stack(trace[:], false)
		}
	}
}
