package config

import (
	"admin_base_server/global"
	"admin_base_server/model/area"
	"admin_base_server/model/config"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"os"
	"sync"
)

type ConfigService struct {
	DB *gorm.DB
}

func NewConfigService() *ConfigService {
	return &ConfigService{DB: global.GVA_DB}
}

func (s *ConfigService) GetActiveConfigFromDB(name string) (map[string]any, error) {
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

func (s *ConfigService) GetQuestionCate() (*config.QuestionCate, error) {
	var (
		cates *config.QuestionCate
		bc    *config.Config
	)
	if err := global.GVA_DB.Where("name = ? AND status = ?", "question_cate", 1).First(&bc).Error; err != nil {
		return nil, err
	}

	if bc == nil || string(bc.Attr) == "" {
		return nil, errors.New("未找到问题类型相关的配置")
	}

	if err := json.Unmarshal(bc.Attr, &cates); err != nil {
		return nil, err
	}

	return cates, nil
}

func (s *ConfigService) GetQuestionLevel() (*config.QuestionLevel, error) {
	var (
		levels *config.QuestionLevel
		bc     *config.Config
	)
	if err := global.GVA_DB.Where("name = ? AND status = ?", "question_level", 1).First(&bc).Error; err != nil {
		return nil, err
	}

	if bc == nil || string(bc.Attr) == "" {
		return nil, errors.New("未找到问题难度相关的配置")
	}

	if err := json.Unmarshal(bc.Attr, &levels); err != nil {
		return nil, err
	}

	return levels, nil
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

func (s *ConfigService) GetAllDBConfigList() (attrs []*config.RConfig, err error) {
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

func (s *ConfigService) GetAreaConfig(level, parentId string) (ret []interface{}, err error) {
	switch level {
	case "province":
		provinces, err := loadProvinces()
		if err != nil {
			return nil, err
		}
		for _, province := range provinces {
			ret = append(ret, province)
		}
	case "city":
		cities, err := loadCities()
		if err != nil {
			return nil, err
		}
		if citiesList, ok := cities[parentId]; ok {
			for _, city := range citiesList {
				ret = append(ret, city)
			}
		} else {
			return nil, errors.New("Invalid parentId parameter")
		}
	case "county":
		counties, err := loadCounties()
		if err != nil {
			return nil, err
		}
		if countiesList, ok := counties[parentId]; ok {
			for _, county := range countiesList {
				ret = append(ret, county)
			}
		} else {
			return nil, errors.New("Invalid level parameter")
		}
	default:
		return nil, errors.New("Invalid parentId parameter")
	}
	return
}

// todo 加上once
func loadProvinces() ([]area.Province, error) {
	data, err := os.ReadFile("./conf_file/area/province.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read province.json: %v", err)
	}
	var p []area.Province
	err = json.Unmarshal(data, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal province.json: %v", err)
	}
	return p, nil
}

func (s *ConfigService) LoadProvincesToMap() (map[string]string, error) {
	provinces, _ := loadProvinces()

	provinceMap := make(map[string]string)
	for _, province := range provinces {
		provinceMap[province.ID] = province.Name
	}

	return provinceMap, nil
}

func loadCities() (map[string][]area.City, error) {
	data, err := os.ReadFile("./conf_file/area/city.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read city.json: %v", err)
	}
	var c map[string][]area.City
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal city.json: %v", err)
	}
	return c, nil
}

var once sync.Once

func (s *ConfigService) LoadCitiesToMap(provinceID string) (map[string]string, error) {
	var cityMap map[string][]area.City
	var err error

	cityMap, err = loadCities()
	if err != nil {
		fmt.Printf("Failed to load cities: %v\n", err)
		return nil, fmt.Errorf("Failed to load cities: %v\n", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load cities: %v", err)
	}

	cities, exists := cityMap[provinceID]
	if !exists {
		return nil, fmt.Errorf("no cities found for province ID %s", provinceID)
	}

	citiesMap := make(map[string]string)
	for _, city := range cities {
		citiesMap[city.ID] = city.Name
	}

	return citiesMap, nil
}

func loadCounties() (map[string][]area.County, error) {
	data, err := os.ReadFile("./conf_file/area/county.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read county.json: %v", err)
	}
	var c map[string][]area.County
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal county.json: %v", err)
	}
	return c, nil
}
