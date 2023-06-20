# cronjob

cronjob is a Go package that allows you to schedule and execute tasks based on cron expressions. It uses Redis for task locking and MySQL for storing task configuration.

## Installation

To use cronjob in your Go project, you can install it using the following command:
```go
    go get github.com/your-username/cronmanager
```

## Usage

Here's an example of how to use cronjob:

```go
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
	//claim redis client
	redisCli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	mysqlCli := mysql.NewClient()
	Start(ctx, redisCli, mysqlCli)
}

```
./example/cron/base.go
```go
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

```
./example/cron/updateInfo.go
```go
import (
	"context"
	"fmt"
	"time"
)

func UpdateInfo(ctx context.Context, _ map[string]interface{}) error {
        //Define task detail
	fmt.Println("hello world")
	fmt.Println("do something......")
	time.Sleep(time.Second * 3)
	fmt.Println("well done")
	return nil
}

```

## Contributing

If you find a bug or want to contribute to CronManager, feel free to open an issue or submit a pull request on GitHub.

## License

CronManager is licensed under the MIT License. See LICENSE for more information.
