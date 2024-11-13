package major

import (
	"gorm.io/gorm"
	"time"
)

// Major 表示专业信息
type Major struct {
	ID                  int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FirstLevelCategory  string    `gorm:"type:varchar(255)" json:"first_level_category"`
	SecondLevelCategory string    `gorm:"type:varchar(255)" json:"second_level_category"`
	MajorName           string    `gorm:"type:varchar(255)" json:"major_name"`
	Year                string    `gorm:"type:varchar(255)" json:"year"`
	CreateTime          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime          time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 返回表名
func (Major) TableName() string {
	return "major"
}

// AutoMigrateMajor 自动迁移表结构
func AutoMigrateMajor(db *gorm.DB) error {
	return db.AutoMigrate(&Major{})
}
