package book

import (
	"github.com/go-playground/validator/v10"
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
	Name            string       `json:"name" validate:"required,max=256"`
	MajorID         int          `json:"major_id" validate:"required,gte=1"`
	MajorName       string       `json:"major_name" validate:"max=256"`
	Level           string       `json:"level" validate:"required,max=100"`
	LevelName       string       `json:"level_name" validate:"max=100"`
	Component       []*Component `json:"component" validate:"required"`
	ComponentDesc   []string     `json:"component_desc"`
	UnitNumber      int          `json:"unit_number" validate:"required,gte=1"`
	QuestionsNumber int          `json:"questions_number"`
	Creator         string       `json:"creator" validate:"required,max=100"`
	Tag             string       `json:"tag" validate:"max=512"`
	CreateTime      time.Time    `json:"create_time"`
	UpdateTime      time.Time    `json:"update_time"`
}

// Validate 验证 RTemplate 结构体
func (rb *RTemplate) Validate() error {
	validate := validator.New()
	return validate.Struct(rb)
}

func (Template) TableName() string {
	return "template"
}

func AutoMigrateTemplate(db *gorm.DB) error {
	return db.AutoMigrate(&Template{})
}
