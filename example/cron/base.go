package cron

import (
	"context"
	"cronjob"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

var TaskMap = map[string]*cronjob.TaskConfInfo{
	"UpdateInfo": {
		CallBack: UpdateInfo,
	},
}

var MisCronManager *cronjob.CronManager

func InitCronTask(ctx context.Context, redisCli *redis.Client, mysqlCli *sql.DB) error {
	defer cronjob.GoRecover(ctx, "Init Cron Task Panic", true)

	MisCronManager = cronjob.NewCronManager(ctx, TaskMap, redisCli, mysqlCli)
	asyncParam := cronjob.AsyncLoadConfigParam{
		Interval: 2,
		CallBack: MisCronManager.LoadConfig,
	}
	go cronjob.AsyncLoadConfig(ctx, &asyncParam)
	go MisCronManager.Run()
	return nil
}
