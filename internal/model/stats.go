package model

import "time"

type Stats struct {
	ID          uint      `gorm:"primaryKey"`                // 主键
	LoginCount  int       `gorm:"not null;default:0"`        // 每天的登录人数
	NewAccounts int       `gorm:"not null;default:0"`        // 每天新增的账号数量
	VideoViews  int       `gorm:"not null;default:0"`        // 每天新增的视频浏览量
	Date        time.Time `gorm:"unique;not null;type:date"` // 日期，确保每天的数据唯一
}
