package warehouse

import (
	"gorm.io/gorm"
	"time"
)

type WarehouseGroup struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Cate        int       `gorm:"not null" json:"cate"`
	WarehouseID string    `gorm:"type:json;not null" json:"warehouse_id"`
	User        string    `gorm:"type:varchar(100);not null" json:"user"`
	CreateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime  time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP" json:"update_time"`
}

func (WarehouseGroup) TableName() string {
	return "warehouse_group"
}

func AutoMigrateWarehouseGroup(db *gorm.DB) error {
	return db.AutoMigrate(&WarehouseGroup{})
}
