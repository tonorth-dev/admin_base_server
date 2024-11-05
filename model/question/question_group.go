package question

import (
	"gorm.io/gorm"
	"time"
)

type QuestionGroup struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Cate       int       `gorm:"not null" json:"cate"`
	QuestionID string    `gorm:"type:json;not null" json:"question_id"`
	User       string    `gorm:"type:varchar(100);not null" json:"user"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

func (QuestionGroup) TableName() string {
	return "question_group"
}

func AutoMigrateQuestionGroup(db *gorm.DB) error {
	return db.AutoMigrate(&QuestionGroup{})
}
