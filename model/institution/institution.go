package institution

import (
	stable "admin_base_server/model/const"
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
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// RInstitution 用于返回给前端的结构体
type RInstitution struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Province     string    `gorm:"type:varchar(255);not null" json:"province"`
	ProvinceName string    `gorm:"type:varchar(255);not null" json:"province_name"`
	City         string    `gorm:"type:varchar(255);not null" json:"city"`
	CityName     string    `gorm:"type:varchar(255);not null" json:"city_name"`
	Password     string    `gorm:"type:varchar(255);not null" json:"password"`
	Leader       string    `gorm:"type:varchar(255);not null" json:"leader"`
	Status       int       `gorm:"not null" json:"status"`
	StatusName   string    `gorm:"not null;default:0" json:"status_name"`
	CreateTime   time.Time `gorm:"not null" json:"create_time"`
	UpdateTime   time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

var StatusMap = map[int]string{
	stable.StatusInitial:  "初始状态",
	stable.StatusDraft:    "未生效",
	stable.StatusActive:   "生效中",
	stable.StatusDeleted:  "已删除",
	stable.StatusAuditing: "审核中",
	stable.StatusExpired:  "已过期",
}

// TableName 返回表名
func (Institution) TableName() string {
	return "institution"
}

// AutoMigrateInstitution 自动迁移表结构
func AutoMigrateInstitution(db *gorm.DB) error {
	return db.AutoMigrate(&Institution{})
}
