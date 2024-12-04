package model

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"updated_at"`
	Username  string `gorm:"size:30;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Salt      string `gorm:"size:255;not null" json:"-"`
	Avatar    string `gorm:"size:255" json:"avatar"`
	Email     string `gorm:"size:255;unique" json:"email"`
	Status    int8   `gorm:"type:tinyint(1);default:0" json:"status"` // 0-正常 1-封禁 2-注销
	Role      int8   `gorm:"type:tinyint(1);default:0" json:"role"`
}
