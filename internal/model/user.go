package model

// User 用户表
type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;<-:create" json:"id"` // 用户 ID
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`       // 创建时间
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"updated_at"`       // 更新时间
	Username  string `gorm:"size:30;not null;unique" json:"username"`      // 用户名
	Password  string `gorm:"size:255;not null" json:"-"`                   // 密码
	Salt      string `gorm:"size:255;not null" json:"-"`                   // 密码盐
	Avatar    string `gorm:"size:255" json:"avatar"`                       // 头像地址
	Email     string `gorm:"size:255;unique" json:"email"`                 // 邮箱
	Status    int8   `gorm:"type:tinyint(1);default:0" json:"status"`      // 用户状态: 0-正常 1-封禁 2-注销
	Role      int8   `gorm:"type:tinyint(1);default:0" json:"role"`        // 用户角色: 0-普通用户 1-管理员
}
