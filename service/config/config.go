package topic

import (
	"admin_base_server/global"
	"admin_base_server/model/config"
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

type ConfigService struct {
	DB *gorm.DB
}

func NewConfigService() *ConfigService {
	return &ConfigService{DB: global.GVA_DB}
}

func (s *ConfigService) GetActiveConfig(name string) (map[string]any, error) {
	var (
		bc   config.Config
		attr map[string]any
	)
	if err := global.GVA_DB.Where("name = ? AND status = ?", name, 1).First(&bc).Error; err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bc.Attr, &attr); err != nil {
		return nil, err
	}

	return attr, nil
}

func (s *ConfigService) UpdateConfig(name string, newAttr map[string]any) error {
	// 将旧的配置状态设置为已失效
	if err := global.GVA_DB.Model(&config.Config{}).Where("name = ? AND status = ?", name, 1).Update("status", 2).Error; err != nil {
		return err
	}

	rawAttr, _ := json.Marshal(newAttr)
	// 插入新的有效配置
	newBC := &config.Config{
		Name:   name,
		Attr:   rawAttr,
		Status: 1,
	}
	return global.GVA_DB.Create(newBC).Error
}

func (s *ConfigService) InsertConfig(name string, attr map[string]any) error {
	// 检查是否有相同的配置名称已存在
	var existingConfig config.Config
	if err := global.GVA_DB.Where("name = ? AND status = ?", name, 1).First(&existingConfig).Error; err == nil {
		// 如果找到相同的配置名称，返回错误
		return errors.New("配置名称已存在")
	}

	// 将 attr 转换为 JSON 格式
	rawAttr, _ := json.Marshal(attr)

	// 创建新的配置
	newConfig := &config.Config{
		Name:   name,
		Attr:   rawAttr,
		Status: 1,
	}

	// 插入新配置
	return global.GVA_DB.Create(newConfig).Error
}

func (s *ConfigService) GetAllConfigList() (attrs []*config.RConfig, err error) {
	var configs []config.Config
	if err := global.GVA_DB.Where("status = ?", 1).Find(&configs).Error; err != nil {
		return nil, err
	}

	for _, cf := range configs {
		var attr map[string]any
		if err := json.Unmarshal(cf.Attr, &attr); err != nil {
			return nil, err
		}
		attrs = append(attrs, &config.RConfig{
			Name:   cf.Name,
			Attr:   attr,
			Status: cf.Status,
		})
	}
	return attrs, nil
}
