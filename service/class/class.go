package class

import (
	"admin_base_server/global"
	"admin_base_server/model/class"
	stable "admin_base_server/model/const"
	"admin_base_server/service/config"
	"admin_base_server/service/institution"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type ClassService struct {
	DB                 *gorm.DB
	configService      *config.ConfigService
	institutionService *institution.InstitutionService
}

func NewClassService() *ClassService {
	return &ClassService{
		DB:                 global.GVA_DB,
		configService:      config.NewConfigService(),
		institutionService: institution.NewInstitutionService(),
	}
}

func (s *ClassService) CreateClass(q *class.Class) error {
	q.Status = stable.StatusActive
	return s.DB.Create(q).Error
}

func (s *ClassService) GetClassByID(id int) (*class.RClass, error) {
	var (
		q *class.Class
		r *class.RClass
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}

	r = &class.RClass{
		ID:            q.ID,
		ClassName:     q.ClassName,
		Password:      q.Password,
		InstitutionID: q.InstitutionID,
		Teacher:       q.Teacher,
		Status:        q.Status,
		StatusName:    stable.RecordStatusMap[q.Status], // 根据实际情况设置 StatusName
		CreateTime:    q.CreateTime,
		UpdateTime:    q.UpdateTime,
	}
	return r, nil
}

func (s *ClassService) GetClassList(page, pageSize int, keyword string, institutionID int) ([]class.RClass, int64, error) {
	var (
		total    int64
		classes  []class.Class
		rClasses []class.RClass
	)

	db := global.GVA_DB.Model(&class.Class{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(class_name) LIKE ? OR LOWER(teacher) LIKE ?", searchQuery, searchQuery)
	}

	// 筛选条件
	if institutionID != 0 {
		db = db.Where("institution_id = ?", institutionID)
	}
	db = db.Where("status = ?", stable.StatusActive)

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range classes {
		insDetail, err := s.institutionService.GetInstitutionByID(v.InstitutionID)
		if err != nil {
			return nil, 0, err
		}

		rClasses = append(rClasses, class.RClass{
			ID:              v.ID,
			ClassName:       v.ClassName,
			Password:        v.Password,
			InstitutionID:   v.InstitutionID,
			InstitutionName: insDetail.Name,
			Teacher:         v.Teacher,
			Status:          v.Status,
			StatusName:      stable.RecordStatusMap[v.Status], // 根据实际情况设置 StatusName
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
		})
	}

	return rClasses, total, nil
}

func (s *ClassService) UpdateClass(id int, q *class.Class) error {
	return s.DB.Model(&class.Class{}).Where("id = ?", id).Updates(q).Error
}

func (s *ClassService) DeleteClass(ids []int) error {
	return s.DB.Model(&class.Class{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
}

// todo 加上校验逻辑
func (s *ClassService) BatchImportClasss(classs []class.Class) error {
	return s.DB.Create(&classs).Error
}

func (s *ClassService) ExportClasss() ([]class.Class, error) {
	var classs []class.Class
	if err := s.DB.Find(&classs).Error; err != nil {
		return nil, err
	}
	return classs, nil
}

func (s *ClassService) getConfigMap() (map[string]string, map[string]string, error) {
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
