package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (ur *User) Create(value *model.User) error {
	return ur.DB.Model(&model.User{}).Create(value).Error
}

func (ur *User) Search(conditions interface{}, limit int, result interface{}) error {
	// return ur.dB.Model(&model.User{}).Where(conditions).Select(fields).Limit(limit).Find(result).Error
	return ur.DB.Model(&model.User{}).Where(conditions).Limit(limit).Find(result).Error
}

func (ur *User) Count(conditions interface{}) (int64, error) {
	var count int64
	err := ur.DB.Model(&model.User{}).Where(conditions).Count(&count).Error
	return count, err
}

func (ur *User) Update(conditions interface{}, fields interface{}, values interface{}) error {
	return ur.DB.Model(&model.User{}).Where(conditions).Select(fields).Updates(values).Error
}

func (ur *User) Delete(conditions interface{}) error {
	return ur.DB.Where(conditions).Delete(&model.User{}).Error
}

// SearchWithOrder 支持排序的分页查询
func (ur *User) SearchWithOrder(conditions interface{}, offset, limit int, order, like string, fields []string, result interface{}) error {
	query := ur.DB.Model(&model.User{})
	if like != "" {
		query = query.Where("username LIKE ?", "%"+like+"%")
	}
	return query.Where(conditions).
		Select(fields).
		Order(order).
		Offset(offset).
		Limit(limit).
		Find(result).Error
}
