package book

import (
	"admin_base_server/model/topic"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type Book struct {
	ID              int         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string      `gorm:"type:varchar(256);not null" json:"name"`
	MajorID         int         `gorm:"not null;default:0" json:"major_id"`
	Level           string      `gorm:"type:varchar(100);not null" json:"level"`
	Component       interface{} `gorm:"type:json;not null" json:"component"`
	UnitNumber      int         `gorm:"not null;default:0" json:"unit_number"`
	Questions       interface{} `gorm:"type:json;not null" json:"questions"`
	QuestionsNumber int         `gorm:"not null;default:0" json:"questions_number"`
	Creator         string      `gorm:"type:varchar(100);not null" json:"creator"`
	TemplateID      int         `gorm:"not null;default:0" json:"template_id"`
	Tag             string      `gorm:"type:varchar(512);not null" json:"tag"`
	CreateTime      time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time   `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RBook struct {
	ID              int                 `json:"id"`
	Name            string              `json:"name" validate:"required,max=256"`
	MajorID         int                 `json:"major_id" validate:"required,gte=1"`
	MajorName       string              `json:"major_name" validate:"max=256"`
	Level           string              `json:"level" validate:"required,max=100"`
	LevelName       string              `json:"level_name" validate:"max=100"`
	Component       []*Component        `json:"component" validate:"required"`
	ComponentDesc   []string            `json:"component_desc"`
	UnitNumber      int                 `json:"unit_number" validate:"required,gte=1"`
	Questions       [][]*Questions      `json:"questions"`
	QuestionsDesc   []*QuestionsDetails `json:"questions_desc"`
	QuestionsNumber int                 `json:"questions_number"`
	Creator         string              `json:"creator" validate:"required,max=100"`
	TemplateID      int                 `json:"template_id"`
	TemplateName    string              `json:"template_name" validate:"max=100"`
	Tag             string              `json:"tag" validate:"max=512"`
	CreateTime      time.Time           `json:"create_time"`
	UpdateTime      time.Time           `json:"update_time"`
}

type Component struct {
	Key    string `json:"key"`
	Number int    `json:"number"`
}

type Questions struct {
	Key string `json:"key"`
	Ids []int  `json:"ids"`
}

type QuestionsDetail struct {
	List     []*topic.RTopic `json:"list"`
	CateName string          `json:"cate_name"`
}

type QuestionsDetails struct {
	Title string             `json:"title"`
	Data  []*QuestionsDetail `json:"questions_detail"`
}

// Validate 验证 RBook 结构体
func (rb *RBook) Validate() error {
	validate := validator.New()
	return validate.Struct(rb)
}

func (Book) TableName() string {
	return "book"
}

func AutoMigrateBook(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}
