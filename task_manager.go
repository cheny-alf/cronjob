package cronjob

import (
	"context"
	"cronjob/logger"
	"cronjob/mysql"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskConfInfo struct {
	Schedue    string //format: * * * * *，represent the M, H, D, M, w, with Unix crontab definition
	Repeatable bool   //True Said whether or not the last time scheduling execution is completed, all new tasks，
	//False Said before the last task, no longer the new task
	Closed   bool                                                //whether the task is shut down,exposed to external can close the cron job
	Params   map[string]interface{}                              //task execution parameter
	CallBack func(context.Context, map[string]interface{}) error //callback function
}

type CronManager struct {
	ctx      context.Context
	redisCli *redis.Client
	mysqlCli *sql.DB
	tasks    sync.Map
}

func NewCronManager(ctx context.Context, tasks map[string]*TaskConfInfo, redisCli *redis.Client, mysqlCli *sql.DB) *CronManager {
	cm := &CronManager{
		ctx:      ctx,
		redisCli: redisCli,
		mysqlCli: mysqlCli,
	}
	if tasks == nil {
		return cm
	}
	for name, info := range tasks {
		cm.tasks.Store(name, info)
	}
	return cm
}

func (cm *CronManager) SetCronTask(taskname string, taskInfo *TaskConfInfo) {
	cm.tasks.Store(taskname, taskInfo)
}

func (cm *CronManager) GetCronTask(taskname string) (*TaskConfInfo, bool) {
	defer GoRecover(cm.ctx, "cron manager get cron task", true)

	info, exist := cm.tasks.Load(taskname)
	if !exist {
		return nil, false
	}
	task := info.(*TaskConfInfo)
	return task, exist
}

func (cm *CronManager) Range(f func(key, value any) bool) {
	cm.tasks.Range(f)
}

const (
	MaxErrorGetPageCount = 5
	GetCronConfigLimit   = 100
)

func (cm *CronManager) LoadConfig(ctx context.Context) {
	defer GoRecover(cm.ctx, "cron manager load config panic", true)

	page, errorGetCount := 0, 0
	for {
		offset := page * GetCronConfigLimit
		taskConfList, err := mysql.GetAllOnlineCronConfigInfo(ctx, offset, GetCronConfigLimit)
		if err != nil {
			logger.Warn(fmt.Sprintf("get all online cron task config faild err:[%v]", err))
			page++
			errorGetCount++
			if errorGetCount >= MaxErrorGetPageCount {
				break
			}
			continue
		}
		if len(taskConfList) == 0 {
			break
		}
		for _, taskConf := range taskConfList {
			task, exist := cm.GetCronTask(taskConf.TaskName)
			if exist {
				if taskConf.CronExpr == "" {
					logger.Warn(fmt.Sprintf("%s task execution cycle not configured, please check cron_expr filed", taskConf.TaskName))
					continue
				}
				param := make(map[string]interface{})
				if taskConf.Params != "" {
					if err := json.Unmarshal([]byte(taskConf.Params), &param); err != nil {
						logger.Warn(fmt.Sprintf("json unmarshal task params error,when load cron task config[%v]", err))
						continue
					}
				}
				newTask := &TaskConfInfo{
					Schedue:    taskConf.CronExpr,
					Repeatable: taskConf.Repeatable == 1,
					Closed:     taskConf.Closed == 1,
					CallBack:   task.CallBack,
					Params:     param,
				}
				cm.SetCronTask(taskConf.TaskName, newTask)
			}
		}
		page++
	}

}

func (cm *CronManager) Run() {
	defer GoRecover(cm.ctx, "cron manager run panic", true)
	cm.LoadConfig(cm.ctx)
	// To perform a task scheduling interval for a second
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			cm.tasks.Range(cm.dispatch)
		case <-cm.ctx.Done():
			return
		}
	}
}

var CronLockKey = "cron_task_lock_%s_%s-%s_%s:%s"

func (cm *CronManager) dispatch(key, value interface{}) bool {
	defer GoRecover(cm.ctx, "cron manager dispatch panic", true)

	curTime := time.Now()

	taskName := key.(string)
	taskConf := value.(*TaskConfInfo)
	// Check whether the current time to execute
	if taskConf.Closed || taskConf.CallBack == nil {
		return true
	}

	if !IsCronTime(curTime, strings.TrimSpace(taskConf.Schedue)) {
		return true
	}
	//  Check the redis locks, judge whether there is a task has been executed
	lockKey := fmt.Sprintf(CronLockKey, taskName, strconv.Itoa(int(curTime.Month())),
		strconv.Itoa(curTime.Day()), strconv.Itoa(curTime.Hour()), strconv.Itoa(curTime.Minute()))
	locker, _ := cm.redisCli.Incr(cm.ctx, lockKey).Result()
	// Lock failed or has been locked
	if locker != 1 {
		return true
	}

	if _, err := cm.redisCli.Expire(cm.ctx, lockKey, time.Minute).Result(); err != nil {
		logger.Warn(fmt.Sprintf("cron task set redis expire error:[%v]", err))
	}

	// do task
	go func(taskName string, task *TaskConfInfo) {
		defer GoRecover(cm.ctx, "do task panic", true)
		taskCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		logger.Info(fmt.Sprintf("Cron_Task_Start:[%s]", taskName))
		errCb := task.CallBack(taskCtx, task.Params)
		if errCb != nil {
			logger.Warn(fmt.Sprintf("cron task set redis expire error:[%v]", errCb))
		}
		logger.Info(fmt.Sprintf("Cron_Task_End:[%s]", taskName))

	}(taskName, taskConf)
	return true
}
