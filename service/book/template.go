package book

import (
	"admin_base_server/global"
	"admin_base_server/model/book"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type TemplateService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
	}
}

func (s *TemplateService) CreateTemplate(q *book.RTemplate) error {
	cateMap, levelMap, err := s.getConfigMap()
	if err != nil {
		return err
	}

	if _, found := levelMap[q.Level]; !found {
		return errors.New("问题等级未定义")
	}

	if len(q.Component) == 0 {
		return errors.New("题本内容不能为空")
	}
	for _, v := range q.Component {
		if _, found := cateMap[v.Key]; !found {
			return errors.New("问题类型未定义")
		}
		if v.Number <= 0 {
			return errors.New("题本内容数量不能小于0")
		}
	}

	componentJSON, _ := json.Marshal(q.Component)

	templateModel := &book.Template{
		ID:         q.ID,
		Name:       q.Name,
		MajorID:    q.MajorID,
		Level:      q.Level,
		Component:  componentJSON,
		UnitNumber: q.UnitNumber,
		Creator:    q.Creator,
	}

	return s.DB.Create(templateModel).Error
}

func (s *TemplateService) GetTemplateByID(id int) (*book.RTemplate, error) {
	var (
		q *book.Template
		r *book.RTemplate
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}
	cateMap, levelMap, err := s.getConfigMap()
	if err != nil {
		return nil, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, err
	}

	var (
		component []*book.Component
	)
	err = json.Unmarshal([]byte(cast.ToString(q.Component)), &component)
	if err != nil {
		return nil, err
	}
	componentDesc, _ := s.generateComponentDesc(component, cateMap)

	r = &book.RTemplate{
		ID:            q.ID,
		Name:          q.Name,
		MajorID:       q.MajorID,
		MajorName:     majorMap[q.MajorID],
		Level:         q.Level,
		LevelName:     levelMap[q.Level],
		Component:     component,
		ComponentDesc: componentDesc,
		UnitNumber:    q.UnitNumber,
		Creator:       q.Creator,
		CreateTime:    q.CreateTime,
		UpdateTime:    q.UpdateTime,
	}
	return r, nil
}

func (s *TemplateService) GetTemplateList(page, pageSize int, keyword, level string, majorID int) ([]book.RTemplate, int64, error) {
	var (
		total      int64
		templates  []book.Template
		rTemplates []book.RTemplate
	)

	db := global.GVA_DB.Model(&book.Template{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(author) LIKE ?", searchQuery, searchQuery)
	}

	// 筛选条件
	if level != "" {
		db = db.Where("level = ?", level)
	}
	if majorID != 0 {
		db = db.Where("major_id = ?", majorID)
	}

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	cateMap, levelMap, err := s.getConfigMap()
	if err != nil {
		return nil, 0, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, 0, err
	}

	for _, v := range templates {
		var (
			component []*book.Component
		)
		err = json.Unmarshal([]byte(cast.ToString(v.Component)), &component)
		if err != nil {
			return nil, 0, err
		}
		componentDesc, _ := s.generateComponentDesc(component, cateMap)

		rTemplates = append(rTemplates, book.RTemplate{
			ID:            v.ID,
			Name:          v.Name,
			MajorID:       v.MajorID,
			MajorName:     majorMap[v.MajorID],
			Level:         v.Level,
			LevelName:     levelMap[v.Level],
			Component:     component,
			ComponentDesc: componentDesc,
			UnitNumber:    v.UnitNumber,
			Creator:       v.Creator,
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		})
	}

	return rTemplates, total, nil
}

func (s *TemplateService) DeleteTemplate(id int) error {
	return s.DB.Delete(&book.Template{}, id).Error
}

func (s *TemplateService) getConfigMap() (map[string]string, map[string]string, error) {
	cateMap := make(map[string]string)
	levelMap := make(map[string]string)

	questionCates, err := s.configService.GetQuestionCate()
	if err != nil {
		return nil, nil, err
	}
	if questionCates == nil || len(questionCates.Cates) == 0 {
		return nil, nil, errors.New("问题类型配置错误")
	}

	questionLevels, err := s.configService.GetQuestionLevel()
	if err != nil {
		return nil, nil, err
	}
	if questionLevels == nil || len(questionLevels.Levels) == 0 {
		return nil, nil, errors.New("问题难度配置错误")
	}

	for _, v := range questionCates.Cates {
		cateMap[v.Id] = v.Name
	}
	for _, v := range questionLevels.Levels {
		levelMap[v.Id] = v.Name
	}

	return cateMap, levelMap, nil
}

func (s *TemplateService) getMajorMap() (map[int]string, error) {
	majorMap := make(map[int]string)

	majors, _, err := s.majorService.GetMajorList(1, 10000, 0, "")
	if err != nil {
		return nil, err
	}

	if len(majors) == 0 {
		return nil, errors.New("专业配置错误")
	}
	for _, v := range majors {
		majorMap[v.ID] = v.MajorName
	}

	return majorMap, nil
}

func (s *TemplateService) generateComponentDesc(component []*book.Component, catMap map[string]string) ([]string, error) {
	var (
		desc []string
	)
	for _, v := range component {
		desc = append(desc, catMap[v.Key]+"：数量"+strconv.Itoa(v.Number))
	}

	return desc, nil
}
