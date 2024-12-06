package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

type Collection struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewCollection(db *gorm.DB) *Collection {
	return &Collection{DB: db}
}

func (r *Collection) Create(value *model.Collection) error {
	return r.DB.Model(&model.Collection{}).Create(value).Error
}

func (r *Collection) Count(conditions interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Collection{}).Where(conditions).Count(&count).Error
	return count, err
}

func (r *Collection) Delete(conditions interface{}) error {
	return r.DB.Where(conditions).Delete(&model.Collection{}).Error
}

func (r *Collection) GetUserCollections(conditions interface{}, limit int, joins []string, fields, result interface{}) error {
	query := r.DB.Model(&model.Collection{}).Where(conditions).Limit(limit).Select(fields)
	for _, join := range joins {
		query = query.Joins(join)
	}
	return query.Find(result).Error
}
