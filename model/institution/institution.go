package institution

import (
	"gorm.io/gorm"
	"time"
)

// Institution 结构体定义
type Institution struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Province   string    `gorm:"type:varchar(255);not null" json:"province"`
	City       string    `gorm:"type:varchar(255);not null" json:"city"`
	Password   string    `gorm:"type:varchar(255);not null" json:"password"`
	Leader     string    `gorm:"type:varchar(255);not null" json:"leader"`
	Status     int       `gorm:"not null" json:"status"`
	CreateTime time.Time `gorm:"not null" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// RInstitution 用于返回给前端的结构体
type RInstitution struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Province   string    `gorm:"type:varchar(255);not null" json:"province"`
	City       string    `gorm:"type:varchar(255);not null" json:"city"`
	Password   string    `gorm:"type:varchar(255);not null" json:"password"`
	Leader     string    `gorm:"type:varchar(255);not null" json:"leader"`
	Status     int       `gorm:"not null" json:"status"`
	StatusName string    `gorm:"not null;default:0" json:"status_name"`
	CreateTime time.Time `gorm:"not null" json:"create_time"`
	UpdateTime time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 返回表名
func (Institution) TableName() string {
	return "institution"
}

// AutoMigrateInstitution 自动迁移表结构
func AutoMigrateInstitution(db *gorm.DB) error {
	return db.AutoMigrate(&Institution{})
}
