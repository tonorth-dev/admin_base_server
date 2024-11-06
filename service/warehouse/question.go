package warehouse

import (
	"admin_base_server/global"
	"admin_base_server/model/warehouse"
	"gorm.io/gorm"
)

type WarehouseService struct {
	DB *gorm.DB
}

func NewWarehouseService() *WarehouseService {
	return &WarehouseService{DB: global.GVA_DB}
}

func (s *WarehouseService) CreateWarehouse(q *warehouse.Warehouse) error {
	return s.DB.Create(q).Error
}

func (s *WarehouseService) GetWarehouseByID(id int) (*warehouse.Warehouse, error) {
	var q warehouse.Warehouse
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}
	return &q, nil
}

func (s *WarehouseService) UpdateWarehouse(id int, q *warehouse.Warehouse) error {
	return s.DB.Model(&warehouse.Warehouse{}).Where("id = ?", id).Updates(q).Error
}

func (s *WarehouseService) DeleteWarehouse(id int) error {
	return s.DB.Delete(&warehouse.Warehouse{}, id).Error
}

func (s *WarehouseService) BatchImportWarehouses(warehouses []warehouse.Warehouse) error {
	return s.DB.Create(&warehouses).Error
}

func (s *WarehouseService) ExportWarehouses() ([]warehouse.Warehouse, error) {
	var warehouses []warehouse.Warehouse
	if err := s.DB.Find(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}
