package cron

import (
	"context"
	"fmt"
	"time"
)

func UpdateInfo(ctx context.Context, _ map[string]interface{}) error {
	fmt.Println("大家好，执行了一次任务")
	fmt.Println("do something......")
	fmt.Println(time.Now())
	return nil
}
