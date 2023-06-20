package mysql

import (
	"context"
	"cronjob/logger"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func NewClient() *sql.DB {
	var err error
	db, err = sql.Open("mysql", "root:128568chen@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	err = db.Ping()
	if err != nil {
		return nil
	}
	return db
}

type CronConfig struct {
	ID           uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"` // 自增主键Id
	TaskName     string `gorm:"column:task_name;NOT NULL"`            // 任务名称, 和代码配置中的相匹配, 用于修改计划任务的调度时间、机房等信息
	CronExpr     string `gorm:"column:cron_expr;NOT NULL"`            // cron的表达式,格式 * * * * *, 分别代表M、H、D、m、w, 同unix crontable定义
	Repeatable   int    `gorm:"column:repeatable;default:0;NOT NULL"` // 任务是否可以重复, 0:不可重复, 1:可以重复
	Closed       int    `gorm:"column:closed;default:0;NOT NULL"`     // 任务是否被关闭, 用于临时关闭任务使用, 0:任务未关闭; 1:任务已关闭
	Params       string `gorm:"column:params;NOT NULL"`               // 任务执行时可能需要的参数
	TaskDescribe string `gorm:"column:task_describe;NOT NULL"`        // 任务描述信息
	Manager      string `gorm:"column:manager;NOT NULL"`              // 任务管理负责人
	IsDeleted    int    `gorm:"column:is_deleted;default:0;NOT NULL"` // 0-正常1-删除
	Ctime        int    `gorm:"column:ctime;NOT NULL"`                // 创建时间
	Mtime        int    `gorm:"column:mtime;NOT NULL"`                // 修改时间
}

func GetAllOnlineCronConfigInfo(ctx context.Context, offset, limit int) ([]CronConfig, error) {
	rows, err := db.Query("SELECT task_name,cron_expr,repeatable,closed,params FROM cron_config LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		logger.Warn(fmt.Sprintf("get all online cron config info failed: %v", err))
	}
	defer rows.Close()

	var configs []CronConfig
	for rows.Next() {
		var c CronConfig
		err = rows.Scan(&c.TaskName, &c.CronExpr, &c.Repeatable, &c.Closed, &c.Params)
		if err != nil {
			logger.Warn(fmt.Sprintf("scan failed: %v", err))
		}
		configs = append(configs, c)
	}
	return configs, err
}
