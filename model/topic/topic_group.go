package topic

import (
	"gorm.io/gorm"
	"time"
)

type TopicGroup struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Cate       int       `gorm:"not null" json:"cate"`
	TopicID    string    `gorm:"type:json;not null" json:"topic_id"`
	User       string    `gorm:"type:varchar(100);not null" json:"user"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

func (TopicGroup) TableName() string {
	return "topic_group"
}

func AutoMigrateTopicGroup(db *gorm.DB) error {
	return db.AutoMigrate(&TopicGroup{})
}
