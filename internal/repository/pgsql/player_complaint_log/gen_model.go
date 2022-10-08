package player_complaint_log

import (
	"time"

	"github.com/lib/pq"
)

// PlayerComplaintLog 玩家投诉记录日日志表
//go:generate gormgen -structs PlayerComplaintLog -input .
type PlayerComplaintLog struct {
	Id                   int            // 商品ID
	LogDate              string         // 添加日期
	GameId               int            // 平台ID
	UserIds              []int          `gorm:"[]int"` // 玩家列表
	Content              string         // 内容
	ImageUrls            pq.StringArray `gorm:"pq.StringArray"`// 图片url列表
	OperateStatus        int            // 操作： 2:封号，1:不处理
	OperateZonstUserId   int            // 最后操作用户ID
	OperateZonstUserName string         //
	OperateTime          time.Time      `gorm:"time"` //
	IsDelete             bool           // 是否已删除
}
