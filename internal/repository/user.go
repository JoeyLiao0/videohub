package repository

import (
	"videohub/internal/model"

	"gorm.io/gorm"
)

// User 提供用户数据访问接口
type User struct {
	DB *gorm.DB
}

// NewUser 实例化用户数据访问对象
func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

// Create 创建用户
func (ur *User) Create(value *model.User) error {
	return ur.DB.Model(&model.User{}).Create(value).Error
}

// Search 查询用户
func (ur *User) Search(conditions interface{}, limit int, result interface{}) error {
	return ur.DB.Model(&model.User{}).Where(conditions).Limit(limit).Find(result).Error
}

// Count 统计用户数量
func (ur *User) Count(conditions interface{}) (int64, error) {
	var count int64
	err := ur.DB.Model(&model.User{}).Where(conditions).Count(&count).Error
	return count, err
}

// Update 更新用户信息
func (ur *User) Update(conditions interface{}, fields interface{}, values interface{}) error {
	return ur.DB.Model(&model.User{}).Where(conditions).Select(fields).Updates(values).Error
}

// Delete 删除用户
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
