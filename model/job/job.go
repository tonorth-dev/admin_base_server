package job

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	ID              int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string          `gorm:"type:varchar(100);not null;default:" json:"code" binding:"required,max=100"`
	Name            string          `gorm:"type:varchar(255);not null;default:" json:"name" binding:"required,max=255"`
	Desc            string          `gorm:"type:varchar(255);not null;default:" json:"desc" binding:"max=255"`
	Cate            string          `gorm:"type:varchar(255);not null;default:" json:"cate" binding:"required,max=255"`
	CompanyCode     string          `gorm:"type:varchar(255);not null;default:" json:"company_code" binding:"required,max=255"`
	CompanyName     string          `gorm:"type:varchar(255);not null;default:" json:"company_name" binding:"required,max=255"`
	EnrollmentNum   int             `gorm:"not null;default:0" json:"enrollment_num" binding:"gte=0"`
	EnrollmentRatio string          `gorm:"type:varchar(32);not null;default:" json:"enrollment_ratio" binding:"max=32"`
	Condition       json.RawMessage `gorm:"type:jsonb" json:"condition" binding:"omitempty"`
	MajorID         int             `gorm:"not null;default:0" json:"major_id" binding:"gte=0"`
	City            string          `gorm:"type:varchar(128);not null;default:" json:"city" binding:"required,max=128"`
	Phone           string          `gorm:"type:varchar(255);not null;default:" json:"phone"`
	Status          int             `gorm:"not null;default:0" json:"status" binding:"gte=0"`
	CreateTime      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time       `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RJob struct {
	ID              int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string          `gorm:"type:varchar(100);not null;default:" json:"code"`
	Name            string          `gorm:"type:varchar(255);not null;default:" json:"name"`
	Desc            string          `gorm:"type:varchar(255);not null;default:" json:"desc"`
	Cate            string          `gorm:"type:varchar(255);not null;default:" json:"cate"`
	CompanyCode     string          `gorm:"type:varchar(255);not null;default:" json:"company_code"`
	CompanyName     string          `gorm:"type:varchar(255);not null;default:" json:"company_name"`
	EnrollmentNum   int             `gorm:"not null;default:0" json:"enrollment_num"`
	EnrollmentRatio string          `gorm:"type:varchar(32);not null;default:" json:"enrollment_ratio"`
	Condition       json.RawMessage `gorm:"type:jsonb" json:"condition"` // 使用 json.RawMessage
	ConditionName   string          `gorm:"type:jsonb" json:"condition_name"`
	MajorID         int             `gorm:"not null;default:0" json:"major_id"`
	MajorName       string          `gorm:"type:varchar(128);not null;default:" json:"major_name"`
	City            string          `gorm:"type:varchar(128);not null;default:" json:"city"`
	Phone           string          `gorm:"type:varchar(255);not null;default:" json:"phone"`
	Status          int             `gorm:"not null;default:0" json:"status"`
	StatusName      string          `gorm:"not null;default:" json:"status_name"`
	CreateTime      time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime      time.Time       `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

var Source = []map[string]string{
	{"id": "1", "name": "高校毕业生"},
	{"id": "2", "name": "社会人才"},
	{"id": "3", "name": "高校毕业生或社会人才"},
}

var Qualification = []map[string]string{
	{"id": "1", "name": "全日制本科以上"},
	{"id": "2", "name": "全日制研究生以上"},
}

var Degree = []map[string]string{
	{"id": "1", "name": "学士以上"},
	{"id": "2", "name": "硕士以上"},
	{"id": "3", "name": "博士以上"},
}

func (Job) TableName() string {
	return "job"
}

func AutoMigrateJob(db *gorm.DB) error {
	return db.AutoMigrate(&Job{})
}
