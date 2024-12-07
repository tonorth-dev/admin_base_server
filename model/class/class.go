package class

import (
	"gorm.io/gorm"
	"time"
)

// Class 结构体定义
type Class struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassName     string    `gorm:"type:varchar(255);not null" json:"class_name"`
	Password      string    `gorm:"type:varchar(255);not null" json:"password"`
	InstitutionID int       `gorm:"not null" json:"institution_id"`
	Teacher       string    `gorm:"type:varchar(255);not null" json:"teacher"`
	Status        bool      `gorm:"not null" json:"status"` // 使用 bool 类型表示 tinyint(1)
	CreateTime    time.Time `gorm:"not null" json:"create_time"`
	UpdateTime    time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// RClass 用于返回给前端的结构体
type RClass struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassName       string    `gorm:"type:varchar(255);not null" json:"class_name"`
	Password        string    `gorm:"type:varchar(255);not null" json:"password"`
	InstitutionID   int       `gorm:"not null" json:"institution_id"`
	InstitutionName int       `gorm:"not null" json:"institution_name"`
	Teacher         string    `gorm:"type:varchar(255);not null" json:"teacher"`
	Status          bool      `gorm:"not null" json:"status"` // 使用 bool 类型表示 tinyint(1)
	StatusName      string    `gorm:"not null;default:0" json:"status_name"`
	CreateTime      time.Time `gorm:"not null" json:"create_time"`
	UpdateTime      time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 返回表名
func (Class) TableName() string {
	return "class"
}

// AutoMigrateClass 自动迁移表结构
func AutoMigrateClass(db *gorm.DB) error {
	return db.AutoMigrate(&Class{})
}
