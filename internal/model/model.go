package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Name     string  `gorm:"size:50;uniqueIndex;not null" json:"name"` // 用户名
	Phone    *string `gorm:"size:20" json:"phone"`                     // 手机号 可为空
	Email    *string `gorm:"size:100;uniqueIndex" json:"email"`        // 邮箱 可为空
	Password string  `gorm:"size:255;not null" json:"-"`               // 密码
	Status   int     `gorm:"default:1" json:"status"`                  // 状态：1-正常，0-禁用
	Type     int     `gorm:"default:1" json:"type"`                    // 类型：1-用户，2-管理员
	Roles    []Role  `gorm:"many2many:account_roles;" json:"roles"`    // 角色
}

type Role struct {
	gorm.Model
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"` // 角色名称
	Description string    `gorm:"size:200" json:"description"`              // 角色描述
	Status      int       `gorm:"default:1" json:"status"`                  // 状态：1-正常，0-禁用
	Accounts    []Account `gorm:"many2many:account_roles;" json:"accounts"`
}
