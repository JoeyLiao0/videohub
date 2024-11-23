package model

type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Salt      string `gorm:"size:255;not null" json:"-"`
	CreatedAt int64  `gorm:"type:bigint;not null;<-:create" json:"created_at"`
	Avatar    string `gorm:"size:255" json:"avatar"`
	Email     string `gorm:"size:255" json:"email"`
	Status    uint8  `gorm:"type:tinyint(1);default:1" json:"status"`
	Role      uint8  `gorm:"type:tinyint(1);default:0" json:"role"`
}
