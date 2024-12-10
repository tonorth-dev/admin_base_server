package student

import (
	"gorm.io/gorm"
	"time"
)

// Student 结构体定义
type Student struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Phone         string    `gorm:"type:varchar(255);not null" json:"phone"`
	Password      string    `gorm:"not null" json:"password"`
	InstitutionID int       `gorm:"not null" json:"institution_id"`
	ClassID       int       `gorm:"not null" json:"class_id"`
	Referrer      string    `gorm:"type:varchar(255);not null" json:"referrer"`
	JobCode       string    `gorm:"type:varchar(255);not null" json:"job_code"`
	MajorIDs      string    `gorm:"type:varchar(512);charset=utf8mb4;collate=utf8mb4_0900_ai_ci;not null" json:"major_ids"`
	Status        int       `gorm:"not null" json:"status"`
	CreateTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

// RStudent 用于返回给前端的结构体
type RStudent struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Phone           string    `gorm:"type:varchar(255);not null" json:"phone"`
	Password        string    `gorm:"not null" json:"password"`
	InstitutionID   int       `gorm:"not null" json:"institution_id"`
	InstitutionName string    `gorm:"not null" json:"institution_name"`
	ClassID         int       `gorm:"not null" json:"class_id"`
	ClassName       string    `gorm:"not null" json:"class_name"`
	Referrer        string    `gorm:"type:varchar(255);not null" json:"referrer"`
	JobCode         string    `gorm:"type:varchar(255);not null" json:"job_code"`
	JobName         string    `gorm:"type:varchar(255);not null" json:"job_name"`
	JobDesc         string    `gorm:"type:varchar(255);not null" json:"job_desc"`
	MajorIDs        []string  `gorm:"type:varchar(512);not null" json:"major_ids"`
	MajorNames      []string  `gorm:"type:varchar(512);not null" json:"major_names"`
	Status          int       `gorm:"not null" json:"status"`
	ClassSorted     int       `gorm:"not null;default:0" json:"class_sorted"`
	StatusName      string    `gorm:"not null;default:0" json:"status_name"`
	CreateTime      time.Time `gorm:"not null" json:"create_time"`
	UpdateTime      time.Time `gorm:"not null;update:CURRENT_TIMESTAMP" json:"update_time"`
}

type RBatchUpdateClass struct {
	StudentIDs []int `json:"student_ids" binding:"required,min=1"`
	ClassId    int   `json:"class_id" binding:"required"`
}

// TableName 返回表名
func (Student) TableName() string {
	return "student"
}

// AutoMigrateStudent 自动迁移表结构
func AutoMigrateStudent(db *gorm.DB) error {
	return db.AutoMigrate(&Student{})
}
