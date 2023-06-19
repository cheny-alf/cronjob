package cron

import (
	"context"
	"fmt"
	"time"
)

func UpdateInfo(ctx context.Context, _ map[string]interface{}) error {
	fmt.Println("hello world")
	fmt.Println("do something......")
	time.Sleep(time.Second * 3)
	fmt.Println("well done")
	return nil
}
