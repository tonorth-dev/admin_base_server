package warehouse

import (
	"admin_base_server/global"
	"admin_base_server/model/warehouse"
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"gorm.io/gorm"
	"os"
)

type WarehouseGroupService struct {
	DB *gorm.DB
}

func NewWarehouseGroupService() *WarehouseGroupService {
	return &WarehouseGroupService{DB: global.GVA_DB}
}

func (s *WarehouseGroupService) CreateWarehouseGroup(g *warehouse.WarehouseGroup) error {
	return s.DB.Create(g).Error
}

func (s *WarehouseGroupService) GetWarehouseGroupByID(id int) (*warehouse.WarehouseGroup, error) {
	var g warehouse.WarehouseGroup
	if err := s.DB.First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *WarehouseGroupService) UpdateWarehouseGroup(id int, g *warehouse.WarehouseGroup) error {
	return s.DB.Model(&warehouse.WarehouseGroup{}).Where("id = ?", id).Updates(g).Error
}

func (s *WarehouseGroupService) DeleteWarehouseGroup(id int) error {
	return s.DB.Delete(&warehouse.WarehouseGroup{}, id).Error
}

func (s *WarehouseGroupService) ExportWarehouseGroupPDF(id int) ([]byte, error) {
	var g warehouse.WarehouseGroup
	if err := s.DB.First(&g, id).Error; err != nil {
		return nil, err
	}

	// Create a temporary blank PDF file
	tmpFile, err := os.CreateTemp("", "warehouse_group_*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	outFile := "/Users/didi/Downloads/warehouse_group.pdf"
	// Use pdfcpu's InsertPage function to add a blank page
	if err := api.InsertPagesFile(tmpFile.Name(), outFile, []string{"1"}, true, nil, api.LoadConfiguration()); err != nil {
		return nil, err
	}

	// Add title as a text watermark
	title := fmt.Sprintf("题组: %s", g.Name)
	titleWM := fmt.Sprintf("f:Helvetica, %s, sc:1 abs, pos:tc, rot:0", title)
	if err := api.AddTextWatermarksFile(outFile, outFile, []string{"1"}, true, titleWM, title, api.LoadConfiguration()); err != nil {
		return nil, err
	}

	// Prepare and add each warehouse as text watermarks
	warehouseIDs := g.WarehouseID
	var ids []int
	if err := json.Unmarshal([]byte(warehouseIDs), &ids); err != nil {
		return nil, err
	}

	for idx, qID := range ids {
		var q warehouse.Warehouse
		if err := s.DB.First(&q, qID).Error; err != nil {
			continue
		}

		// Format warehouse details
		warehouseText := fmt.Sprintf("题目: %s | 类型: %d | 答案: %s | 录入人: %s | 专业ID: %d | 专业名称: %s | 标签: %s",
			q.Title, q.Cate, q.Answer, q.Author, q.MajorID, q.MajorName, q.Tag)

		// Position each warehouse below the previous text
		position := fmt.Sprintf("f:Helvetica, %s, sc:0.8 abs, pos:tl, dy:%d", warehouseText, -100*(idx+1))
		if err := api.AddTextWatermarksFile(outFile, outFile, []string{"1"}, true, titleWM, position, api.LoadConfiguration()); err != nil {
			return nil, err
		}
	}

	// Read the resulting PDF data
	pdfData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	return pdfData, nil
}
