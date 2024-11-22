package topic

import (
	"gorm.io/gorm"
	"time"
)

type Topic struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(10000);not null" json:"title"`
	Cate       string    `gorm:"not null" json:"cate"`
	Level      string    `gorm:"not null" json:"level"`
	Answer     string    `gorm:"type:text;not null" json:"answer"`
	Author     string    `gorm:"type:varchar(100);not null" json:"author"`
	MajorID    int       `gorm:"not null" json:"major_id"`
	Tag        string    `gorm:"type:varchar(255)" json:"tag"`
	Status     int       `gorm:"not null" json:"status" validate:"required,oneof=1 2"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RTopic struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(10000);not null" json:"title"`
	Cate       string    `gorm:"not null" json:"cate"`
	CateName   string    `gorm:"not null" json:"cate_name"`
	Level      string    `gorm:"not null" json:"level"`
	LevelName  string    `gorm:"not null" json:"level_name"`
	Answer     string    `gorm:"type:text;not null" json:"answer"`
	Author     string    `gorm:"type:varchar(100);not null" json:"author"`
	MajorID    int       `gorm:"not null" json:"major_id"`
	MajorName  string    `gorm:"not null" json:"major_name"`
	Tag        string    `gorm:"type:varchar(255)" json:"tag"`
	Status     int       `gorm:"not null" json:"status"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

func (Topic) TableName() string {
	return "topic"
}

func AutoMigrateTopic(db *gorm.DB) error {
	return db.AutoMigrate(&Topic{})
}
