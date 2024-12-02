package student

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/student"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type StudentService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
}

func NewStudentService() *StudentService {
	return &StudentService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
	}
}

func (s *StudentService) CreateStudent(q *student.Student) error {
	return s.DB.Create(q).Error
}

func (s *StudentService) GetStudentByID(id int) (*student.RStudent, error) {
	var (
		q *student.Student
		r *student.RStudent
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}

	r = &student.RStudent{
		ID:            q.ID,
		Name:          q.Name,
		Phone:         q.Phone,
		Password:      q.Password,
		InstitutionID: q.InstitutionID,
		ClassID:       q.ClassID,
		Referrer:      q.Referrer,
		Status:        q.Status,
		StatusName:    "Active", // 根据实际情况设置 StatusName
		CreateTime:    q.CreateTime,
		UpdateTime:    q.UpdateTime,
	}
	return r, nil
}

func (s *StudentService) GetStudentList(page, pageSize int, keyword, cate, level string, majorID int, status int) ([]student.RStudent, int64, error) {
	var (
		total     int64
		students  []student.Student
		rStudents []student.RStudent
	)

	db := global.GVA_DB.Model(&student.Student{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(company_name) LIKE ? OR LOWER(student_name) LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	// 筛选条件
	if level != "" {
		db = db.Where("level = ?", level)
	}
	if cate != "" {
		db = db.Where("cate = ?", cate)
	}
	if majorID != 0 {
		db = db.Where("major_id = ?", majorID)
	}
	db = db.Where("status = ?", stable.StatusActive)

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	for _, v := range students {
		rStudents = append(rStudents, student.RStudent{
			ID:            v.ID,
			Name:          v.Name,
			Phone:         v.Phone,
			Password:      v.Password,
			InstitutionID: v.InstitutionID,
			ClassID:       v.ClassID,
			Referrer:      v.Referrer,
			Status:        v.Status,
			StatusName:    "Active", // 根据实际情况设置 StatusName
			CreateTime:    v.CreateTime,
			UpdateTime:    v.UpdateTime,
		})
	}

	return rStudents, total, nil
}

func (s *StudentService) UpdateStudent(id int, q *student.Student) error {
	return s.DB.Model(&student.Student{}).Where("id = ?", id).Updates(q).Error
}

func (s *StudentService) DeleteStudent(ids []int) error {
	return s.DB.Model(&student.Student{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
}

// todo 加上校验逻辑
func (s *StudentService) BatchImportStudents(students []student.Student) error {
	return s.DB.Create(&students).Error
}

func (s *StudentService) ExportStudents() ([]student.Student, error) {
	var students []student.Student
	if err := s.DB.Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (s *StudentService) getConfigMap() (map[string]string, map[string]string, error) {
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

func (s *StudentService) getMajorMap() (map[int]string, error) {
	majorMap := make(map[int]string)

	majors, _, err := s.majorService.GetMajorList(1, 10000, "")
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
