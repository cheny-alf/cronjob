package main

import (
	"context"
	"cronjob/example/cron"
	"cronjob/logger"
	"cronjob/mysql"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
)

func Start(ctx context.Context, redisCli *redis.Client, mysqlCli *sql.DB) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() error {
		return cron.InitCronTask(ctx, redisCli, mysqlCli)
	}()

	wg.Wait()
	fmt.Println("end")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "cronjob",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	//声明redis client
	redisCli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	mysqlCli := mysql.NewClient()
	Start(ctx, redisCli, mysqlCli)
}
