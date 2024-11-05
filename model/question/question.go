package question

import (
	"gorm.io/gorm"
	"time"
)

type Question struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(10000);not null" json:"title"`
	Cate       int       `gorm:"not null" json:"cate"`
	Answer     string    `gorm:"type:text;not null" json:"answer"`
	Author     string    `gorm:"type:varchar(100);not null" json:"author"`
	MajorID    int       `gorm:"not null" json:"major_id"`
	MajorName  string    `gorm:"type:varchar(100);not null" json:"major_name"`
	Tag        string    `gorm:"type:varchar(255)" json:"tag"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

func (Question) TableName() string {
	return "question"
}

func AutoMigrateQuestion(db *gorm.DB) error {
	return db.AutoMigrate(&QuestionGroup{})
}
