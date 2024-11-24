package book

import (
	"gorm.io/gorm"
	"time"
)

type Template struct {
	ID         int         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string      `gorm:"type:varchar(256);not null" json:"name"`
	MajorID    int         `gorm:"not null;default:0" json:"major_id"`
	Level      string      `gorm:"type:varchar(100);not null" json:"level"`
	Component  interface{} `gorm:"type:json;not null" json:"component"`
	UnitNumber int         `gorm:"not null;default:0" json:"unit_number"`
	Creator    string      `gorm:"type:varchar(100);not null" json:"creator"`
	CreateTime time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time   `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RTemplate struct {
	ID              int          `json:"id"`
	Name            string       `json:"name" binding:"required,max=256"`
	MajorID         int          `json:"major_id" binding:"required,gte=1"`
	MajorName       string       `json:"major_name" binding:"max=256"`
	Level           string       `json:"level" binding:"required,max=100"`
	LevelName       string       `json:"level_name" binding:"max=100"`
	Component       []*Component `json:"component" binding:"required"`
	ComponentDesc   []string     `json:"component_desc"`
	UnitNumber      int          `json:"unit_number" binding:"required,gte=1"`
	QuestionsNumber int          `json:"questions_number"`
	Creator         string       `json:"creator" binding:"required,max=100"`
	Tag             string       `json:"tag" binding:"max=512"`
	CreateTime      time.Time    `json:"create_time"`
	UpdateTime      time.Time    `json:"update_time"`
}

func (Template) TableName() string {
	return "template"
}

func AutoMigrateTemplate(db *gorm.DB) error {
	return db.AutoMigrate(&Template{})
}
