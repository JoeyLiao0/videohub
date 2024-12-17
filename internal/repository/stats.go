package repository

import (
	"time"
	"videohub/global"
	"videohub/internal/model"

	"github.com/redis/go-redis/v9"
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
	if err := global.Rdb.Del(global.Ctx, "login_users").Err(); err != nil {
		logrus.Error(err.Error())
		return
	}
	
	var newAccounts int64
	if err := db.Model(&model.User{}).Where("created_at BETWEEN ? AND ? AND role = 0", start, end).Count(&newAccounts).Error; err != nil {
		logrus.Error(err.Error())
		return
	}

	var currentViews int64
	var cursor uint64
	pattern := "video:*:views"

	for {
		keys, nextCursor, err := global.Rdb.Scan(global.Ctx, cursor, pattern, 100).Result()
		if err != nil {
			logrus.Error(err.Error())
			return
		}

		for _, key := range keys {
			viewCount, err := global.Rdb.Get(global.Ctx, key).Int64()
			if err != nil {
				logrus.Error(err.Error())
				return
			}
			currentViews += viewCount
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	lastViews, err := global.Rdb.Get(global.Ctx, "video:total_views").Int64()
	if err == redis.Nil {
		lastViews = 0
	} else if err != nil {
		logrus.Error(err.Error())
		return
	}

	key := "video:total_views"
	if err := global.Rdb.Set(global.Ctx, key, currentViews, 0).Err(); err != nil {
		logrus.Error(err.Error())
		return
	}

	newViews := currentViews - lastViews
	stat := model.Stats{
		LoginCount:    int(loginCount),
		NewAccounts:   int(newAccounts),
		NewVideoViews: int(newViews),
		Date:          time.Now(),
	}
	if err := db.Model(&model.Stats{}).Create(&stat).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (e *Stats) Search(startDate, endDate string, limit int, result interface{}) error {
	return e.DB.Model(&model.Stats{}).Where("date BETWEEN ? AND ?", startDate, endDate).Limit(limit).Find(result).Error
}
