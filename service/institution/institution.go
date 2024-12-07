package institution

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/institution"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"crypto/rand"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/big"
	"strings"
)

type InstitutionService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
}

func NewInstitutionService() *InstitutionService {
	return &InstitutionService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
	}
}

func (s *InstitutionService) CreateInstitution(q *institution.Institution) error {
	q.Password, _ = s.generateUniquePassword(12)
	return s.DB.Create(q).Error
}

func (s *InstitutionService) GetInstitutionByID(id int) (*institution.RInstitution, error) {
	var (
		q *institution.Institution
		r *institution.RInstitution
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}

	provincesMap, _ := s.configService.LoadProvincesToMap()
	cityMap, _ := s.configService.LoadCitiesToMap(q.Province)

	r = &institution.RInstitution{
		ID:           q.ID,
		Name:         q.Name,
		Province:     q.Province,
		ProvinceName: provincesMap[q.Province],
		City:         q.City,
		CityName:     cityMap[q.City],
		Password:     q.Password,
		Leader:       q.Leader,
		Status:       q.Status,
		StatusName:   institution.StatusMap[q.Status],
		CreateTime:   q.CreateTime,
		UpdateTime:   q.UpdateTime,
	}
	return r, nil
}

func (s *InstitutionService) GetInstitutionList(page, pageSize int, keyword, province, city string) ([]institution.RInstitution, int64, error) {
	var (
		total         int64
		institutions  []institution.Institution
		rInstitutions []institution.RInstitution
	)

	db := global.GVA_DB.Model(&institution.Institution{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(leader) LIKE ?", searchQuery, searchQuery)
	}

	// 筛选条件
	if province != "" {
		db = db.Where("province = ?", province)
	}
	if city != "" {
		db = db.Where("city = ?", city)
	}
	db = db.Where("status != ?", stable.StatusDeleted)
	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&institutions).Error; err != nil {
		return nil, 0, err
	}

	provincesMap, _ := s.configService.LoadProvincesToMap()

	for _, v := range institutions {
		cityMap, _ := s.configService.LoadCitiesToMap(v.Province)
		rInstitutions = append(rInstitutions, institution.RInstitution{
			ID:           v.ID,
			Name:         v.Name,
			Province:     v.Province,
			ProvinceName: provincesMap[v.Province],
			City:         v.City,
			CityName:     cityMap[v.City],
			Password:     v.Password,
			Leader:       v.Leader,
			Status:       v.Status,
			StatusName:   institution.StatusMap[v.Status],
			CreateTime:   v.CreateTime,
			UpdateTime:   v.UpdateTime,
		})
	}

	return rInstitutions, total, nil
}

func (s *InstitutionService) UpdateInstitution(id int, q *institution.Institution) error {
	return s.DB.Model(&institution.Institution{}).Where("id = ?", id).Updates(q).Error
}

func (s *InstitutionService) DeleteInstitution(ids []int) error {
	return s.DB.Model(&institution.Institution{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
}

// todo 加上校验逻辑
func (s *InstitutionService) BatchImportInstitutions(institutions []institution.Institution) error {
	return s.DB.Create(&institutions).Error
}

func (s *InstitutionService) ExportInstitutions() ([]institution.Institution, error) {
	var institutions []institution.Institution
	if err := s.DB.Find(&institutions).Error; err != nil {
		return nil, err
	}
	return institutions, nil
}

func (s *InstitutionService) getConfigMap() (map[string]string, map[string]string, error) {
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

func (s *InstitutionService) getMajorMap() (map[int]string, error) {
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

func (s *InstitutionService) generateUniquePassword(length int) (string, error) {
	// 定义字符集，包括大小写字母、数字和特殊字符
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	charsetLength := len(charset)

	// 确保长度合理
	if length < 4 {
		return "", fmt.Errorf("密码长度不能小于 4")
	}

	// 随机生成密码
	password := make([]byte, length)
	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(charsetLength)))
		if err != nil {
			return "", fmt.Errorf("随机数生成失败: %v", err)
		}
		password[i] = charset[index.Int64()]
	}

	// 返回最终密码
	return string(password), nil
}
