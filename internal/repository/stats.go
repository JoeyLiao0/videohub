package repository

import (
	"time"
	"videohub/global"
	"videohub/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Stats struct {
	DB *gorm.DB
}

func NewStats(db *gorm.DB) *Stats {
	return &Stats{DB: db}
}

func WriteStats(db *gorm.DB) {
	endOfDay := time.Now().Truncate(24 * time.Hour)
	startOfDay := endOfDay.Add(-24 * time.Hour)
	start := startOfDay.UnixMilli()
	end := endOfDay.UnixMilli()

	loginCount, err := global.Rdb.SCard(global.Ctx, "login_users").Result()
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	var newAccounts int64
	if err := db.Model(&model.User{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&newAccounts).Error; err != nil {
		logrus.Error(err.Error())
		return
	}

	var videoViews int64
	if err := db.Model(&model.Video{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&videoViews).Error; err != nil {
		logrus.Error(err.Error())
		return
	}

	stat := model.Stats{
		LoginCount:  int(loginCount),
		NewAccounts: int(newAccounts),
		VideoViews:  int(videoViews),
		Date:        time.Now(),
	}
	if err := db.Model(&model.Stats{}).Create(&stat).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (e *Stats) Search(startDate, endDate string, limit int, result interface{}) error {
	return e.DB.Model(&model.Stats{}).Where("date BETWEEN ? AND ?", startDate, endDate).Limit(limit).Find(result).Error
}
