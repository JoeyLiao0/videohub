package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

type User struct {
	dB *gorm.DB
}

// 工厂函数，存储单例的数据库操作对象
func NewUser(db *gorm.DB) *User {
	return &User{dB: db}
}

func (ur *User) Create(value *model.User) error {
	return ur.dB.Model(&model.User{}).Create(value).Error
}

func (ur *User) Search(conditions interface{}, limit int, result interface{}) error {
	// return ur.dB.Model(&model.User{}).Where(conditions).Select(fields).Limit(limit).Find(result).Error
	return ur.dB.Model(&model.User{}).Where(conditions).Limit(limit).Find(result).Error
}

func (ur *User) Count(conditions interface{}) (int64, error) {
	var count int64
	err := ur.dB.Model(&model.User{}).Where(conditions).Count(&count).Error
	return count, err
}

func (ur *User) Update(conditions interface{}, fields interface{}, values interface{}) error {
	return ur.dB.Model(&model.User{}).Where(conditions).Select(fields).Updates(values).Error
}

func (ur *User) Delete(conditions interface{}) error {
	return ur.dB.Where(conditions).Delete(&model.User{}).Error
}
