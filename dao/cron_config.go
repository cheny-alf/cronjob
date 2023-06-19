package dao

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

func (m CronConfig) TableName() string {
	return "cron_config"
}

//func GetAllOnlineCronConfigInfo(ctx context.Context, offset, limit int, filed []string) ([]CronConfig, error) {
//	//考虑limit 限制避免一次数据返回过多
//	obj := new(CronConfig)
//	db := mysql.GetTable(ctx, obj)
//	var result []CronConfig
//	err := db.Select(filed).Where("is_deleted = 0").Offset(offset).Limit(limit).Find(&result).Error
//	return result, err
//}
