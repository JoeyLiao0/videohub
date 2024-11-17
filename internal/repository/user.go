package repository

import (
	"fmt"
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

func (ur *User) SearchByUsername(username string, user *model.User) error {
	return ur.dB.Where("username = ?", username).First(user).Error
}

func (ur *User) SearchById(id uint64, user *model.User) error {
	return ur.dB.Where("id = ?", id).First(user).Error
}

// GetUserByID 根据用户 ID 查找用户信息
func (ur *User) GetUserByID(userID uint64) (*model.User, error) {
	var user model.User
	// 使用 GORM 的 First() 方法，根据用户 ID 查找用户信息
	if err := ur.dB.First(&user, userID).Error; err != nil {
		// 如果未找到用户，返回自定义错误信息
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户未找到")
		}
		// 返回其他错误
		return nil, err
	}
	// 返回找到的用户信息
	return &user, nil
}

// CreateUser 创建新用户
func (ur *User) CreateUser(newUser *model.User) error {
	// 使用 GORM 的 Create() 方法将新用户插入数据库
	if err := ur.dB.Create(newUser).Error; err != nil {
		// 如果创建用户时出现错误，则返回错误
		return err
	}
	// 如果成功，则返回 nil（无错误）
	return nil
}

// UpdateUser 更新用户信息
func (ur *User) UpdateUser(updatedUser *model.User) error {
	// 使用 GORM 的 Save() 方法更新用户信息
	if err := ur.dB.Save(updatedUser).Error; err != nil {
		// 如果更新用户信息时出现错误，则返回错误
		return err
	}
	// 如果更新成功，则返回 nil（无错误）
	return nil
}

// DeleteUser 删除用户，根据用户 ID
func (ur *User) DeleteUser(userID uint) error {
	// 使用 GORM 的 Delete() 方法根据用户 ID 删除用户
	if err := ur.dB.Delete(&model.User{}, userID).Error; err != nil {
		// 如果删除用户时出现错误，则返回错误
		return err
	}
	// 如果删除成功，则返回 nil（无错误）
	return nil
}

// FindUserByUsername 根据用户名查找用户
func (ur *User) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	// 使用 GORM 的 Where() 方法根据用户名查找用户信息
	if err := ur.dB.Where("username = ?", username).First(&user).Error; err != nil {
		// 如果未找到用户，返回 nil 而不返回错误（用户不存在不是错误）
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// 返回其他错误
		return nil, err
	}
	// 返回找到的用户信息
	return &user, nil
}

// FindUserByEmail 根据邮箱查找用户
func (ur *User) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	// 使用 GORM 的 Where() 方法根据邮箱查找用户信息
	if err := ur.dB.Where("email = ?", email).First(&user).Error; err != nil {
		// 如果未找到用户，返回 nil 而不返回错误（用户不存在不是错误）
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// 返回其他错误
		return nil, err
	}
	// 返回找到的用户信息
	return &user, nil
}

// FindAllUsers 获取所有用户信息
func (ur *User) GetAllUsers() ([]model.User, error) {
	var users []model.User
	// 使用 GORM 的 Find() 方法查找所有用户信息
	if err := ur.dB.Find(&users).Error; err != nil {
		// 如果查找用户时出现错误，则返回错误
		return nil, err
	}
	// 返回用户列表
	return users, nil
}