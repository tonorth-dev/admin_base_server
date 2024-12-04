package major

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/major"
	"gorm.io/gorm"
	"strings"
)

// MajorService 提供对 Major 模型的操作
type MajorService struct {
	DB *gorm.DB
}

// NewMajorService 创建一个新的 MajorService 实例
func NewMajorService() *MajorService {
	return &MajorService{DB: global.GVA_DB}
}

// CreateMajor 创建新的专业信息
func (s *MajorService) CreateMajor(m *major.Major) error {
	m.Status = stable.StatusActive
	return s.DB.Create(m).Error
}

// GetMajorByID 根据 ID 获取专业信息
func (s *MajorService) GetMajorByID(id int) (*major.Major, error) {
	var m major.Major
	if err := s.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMajorList 获取专业信息列表
func (s *MajorService) GetMajorList(page, pageSize int, id int, keyword string) ([]major.Major, int64, error) {
	var majors []major.Major
	var total int64

	db := global.GVA_DB.Model(&major.Major{})

	if id > 0 {
		db = db.Where("id = ?", id)
	}

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(first_level_category) LIKE ? OR LOWER(second_level_category) LIKE ? OR LOWER(major_name) LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&majors).Error; err != nil {
		return nil, 0, err
	}

	return majors, total, nil
}

// UpdateMajor 更新专业信息
func (s *MajorService) UpdateMajor(id int, m *major.Major) error {
	return s.DB.Model(&major.Major{}).Where("id = ?", id).Updates(m).Error
}

// DeleteMajor 删除专业信息
func (s *MajorService) DeleteMajor(id int) error {
	return s.DB.Delete(&major.Major{}, id).Error
}

// BatchImportMajors 批量导入专业信息
func (s *MajorService) BatchImportMajors(majors []major.Major) error {
	return s.DB.Create(&majors).Error
}

// ExportMajors 导出专业信息
func (s *MajorService) ExportMajors() ([]major.Major, error) {
	var majors []major.Major
	if err := s.DB.Find(&majors).Error; err != nil {
		return nil, err
	}
	return majors, nil
}
