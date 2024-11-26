package topic

import (
	"gorm.io/gorm"
	"time"
)

type Topic struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(10000);not null" json:"title" binding:"required,max=10000"`
	Cate        string    `gorm:"not null" json:"cate" binding:"required,max=255"`
	Level       string    `gorm:"not null" json:"level" binding:"required,max=255"`
	Answer      string    `gorm:"type:text;not null" json:"answer"`
	AnswerDraft string    `gorm:"type:text;not null" json:"answer_draft"`
	Author      string    `gorm:"type:varchar(100);not null" json:"author" binding:"required,max=100"`
	Invitee     string    `gorm:"type:varchar(100);not null" json:"invitee"`
	MajorID     int       `gorm:"not null" json:"major_id" binding:"required,gt=0"`
	Tag         string    `gorm:"type:varchar(255)" json:"tag" binding:"max=255"`
	Status      int       `gorm:"not null" json:"status" binding:"required,oneof=1 2"`
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RTopic struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(10000);not null" json:"title"`
	Cate        string    `gorm:"not null" json:"cate"`
	CateName    string    `gorm:"not null" json:"cate_name"`
	Level       string    `gorm:"not null" json:"level"`
	LevelName   string    `gorm:"not null" json:"level_name"`
	Answer      string    `gorm:"type:text;not null" json:"answer"`
	AnswerDraft string    `gorm:"type:text;not null" json:"answer_draft"`
	Author      string    `gorm:"type:varchar(100);not null" json:"author"`
	Invitee     string    `gorm:"type:varchar(100);not null" json:"invitee"`
	MajorID     int       `gorm:"not null" json:"major_id"`
	MajorName   string    `gorm:"not null" json:"major_name"`
	Tag         string    `gorm:"type:varchar(255)" json:"tag"`
	Status      int       `gorm:"not null" json:"status"`
	StatusName  string    `gorm:"not null" json:"status_name"`
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type Audit struct {
	Answer  string `gorm:"type:text;not null" json:"answer" binding:"required,max=10000"`
	Invitee string `gorm:"type:varchar(100);not null" json:"invitee" binding:"required,max=100"`
}

func (Topic) TableName() string {
	return "topic"
}

func AutoMigrateTopic(db *gorm.DB) error {
	return db.AutoMigrate(&Topic{})
}
