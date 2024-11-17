package config

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	ID         uint            `gorm:"primaryKey" json:"id"`
	Name       string          `gorm:"uniqueIndex;not null" json:"name"`
	Attr       json.RawMessage `gorm:"type:jsonb" json:"attr"` // 使用 json.RawMessage
	Status     int             `gorm:"default:0" json:"status"`
	CreateTime time.Time       `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time       `gorm:"autoUpdateTime" json:"update_time"`
}

type RConfig struct {
	Name   string         `gorm:"uniqueIndex;not null" json:"name"`
	Attr   map[string]any `gorm:"type:jsonb" json:"attr"` // 使用 json.RawMessage
	Status int            `gorm:"default:true" json:"status"`
}

type QuestionCate struct {
	Cates []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"cates"`
}

type QuestionLevel struct {
	Levels []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"levels"`
}

// TableName 返回表名
func (Config) TableName() string {
	return "business_config"
}

// AutoMigrateMajor 自动迁移表结构
func AutoMigrateMajor(db *gorm.DB) error {
	return db.AutoMigrate(&Config{})
}
