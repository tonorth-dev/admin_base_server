package student

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/student"
	"admin_base_server/service/config"
	"admin_base_server/service/institution"
	"admin_base_server/service/job"
	"admin_base_server/service/major"
	"errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"log"
	"strings"
)

type StudentService struct {
	DB                 *gorm.DB
	configService      *config.ConfigService
	majorService       *major.MajorService
	jobService         *job.JobService
	institutionService *institution.InstitutionService
}

func NewStudentService() *StudentService {
	return &StudentService{
		DB:                 global.GVA_DB,
		configService:      config.NewConfigService(),
		majorService:       major.NewMajorService(),
		jobService:         job.NewJobService(),
		institutionService: institution.NewInstitutionService(),
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

	// 获取 JobName 和 JobDesc
	jobDetail, err := s.jobService.GetJobByCode(q.JobCode)
	if err != nil {
		return nil, err
	}
	insDetail, err := s.institutionService.GetInstitutionByID(q.InstitutionID)
	if err != nil {
		return nil, err
	}

	// 获取 MajorNames
	majorMap, _ := s.getMajorMap()
	majorIDs := strings.Split(q.MajorIDs, ",")
	var majorNames []string

	// 遍历 majorIDs 并使用 majorMap 获取对应的 majorName
	for _, majorID := range majorIDs {
		if majorName, exists := majorMap[cast.ToInt(majorID)]; exists {
			majorNames = append(majorNames, majorName)
		} else {
			// 如果 majorID 在 majorMap 中不存在，可以选择记录日志或跳过
			log.Printf("MajorID %s 在 majorMap 中不存在", majorID)
		}
	}

	r = &student.RStudent{
		ID:              q.ID,
		Name:            q.Name,
		Phone:           q.Phone,
		Password:        q.Password,
		InstitutionID:   q.InstitutionID,
		InstitutionName: insDetail.Name,
		ClassID:         q.ClassID,
		ClassName:       "s.getClassName(q.ClassID)",
		Referrer:        q.Referrer,
		JobCode:         q.JobCode,
		JobName:         jobDetail.Name,
		JobDesc:         jobDetail.Desc,
		MajorIDs:        q.MajorIDs,
		MajorNames:      majorNames,
		Status:          q.Status,
		StatusName:      stable.RecordStatusMap[q.Status],
		CreateTime:      q.CreateTime,
		UpdateTime:      q.UpdateTime,
	}
	return r, nil
}

func (s *StudentService) GetStudentList(page, pageSize int, keyword string, classId, institutionId, majorID int, status int) ([]student.RStudent, int64, error) {
	var (
		total     int64
		students  []student.Student
		rStudents []student.RStudent
	)

	db := global.GVA_DB.Model(&student.Student{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(job_code) LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	// 筛选条件
	if classId != 0 {
		db = db.Where("class_id = ?", classId)
	}
	if institutionId != 0 {
		db = db.Where("institution_id = ?", institutionId)
	}
	if majorID != 0 {
		db = db.Where("FIND_IN_SET(?, major_ids)", majorID)
	}
	if status != 0 {
		db = db.Where("status =?", status)
	} else {
		db = db.Where("status != ?", stable.StatusDeleted)
	}

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	majorMap, _ := s.getMajorMap()
	for _, v := range students {

		// 获取 JobName 和 JobDesc
		jobDetail, err := s.jobService.GetJobByCode(v.JobCode)
		if err != nil {
			return nil, 0, err
		}
		insDetail, err := s.institutionService.GetInstitutionByID(v.InstitutionID)
		if err != nil {
			return nil, 0, err
		}

		// 获取 MajorNames
		majorIDs := strings.Split(v.MajorIDs, ",")
		var majorNames []string

		// 遍历 majorIDs 并使用 majorMap 获取对应的 majorName
		for _, majorID := range majorIDs {
			if majorName, exists := majorMap[cast.ToInt(majorID)]; exists {
				majorNames = append(majorNames, majorName)
			} else {
				// 如果 majorID 在 majorMap 中不存在，可以选择记录日志或跳过
				log.Printf("MajorID %s 在 majorMap 中不存在", majorID)
			}
		}

		rStudents = append(rStudents, student.RStudent{
			ID:              v.ID,
			Name:            v.Name,
			Phone:           v.Phone,
			Password:        v.Password,
			InstitutionID:   v.InstitutionID,
			InstitutionName: insDetail.Name,
			ClassID:         v.ClassID,
			ClassName:       "s.getClassName(v.ClassID)",
			Referrer:        v.Referrer,
			JobCode:         v.JobCode,
			JobName:         jobDetail.Name,
			JobDesc:         jobDetail.Desc,
			MajorIDs:        v.MajorIDs,
			MajorNames:      majorNames,
			Status:          v.Status,
			StatusName:      stable.RecordStatusMap[v.Status],
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
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
