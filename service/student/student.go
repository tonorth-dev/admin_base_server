package student

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/student"
	"admin_base_server/service/class"
	"admin_base_server/service/config"
	"admin_base_server/service/institution"
	"admin_base_server/service/job"
	"admin_base_server/service/major"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"log"
	"math/big"
	"strings"
)

type StudentService struct {
	DB                 *gorm.DB
	configService      *config.ConfigService
	majorService       *major.MajorService
	jobService         *job.JobService
	institutionService *institution.InstitutionService
	classService       *class.ClassService
}

func NewStudentService() *StudentService {
	return &StudentService{
		DB:                 global.GVA_DB,
		configService:      config.NewConfigService(),
		majorService:       major.NewMajorService(),
		jobService:         job.NewJobService(),
		institutionService: institution.NewInstitutionService(),
		classService:       class.NewClassService(),
	}
}

func (s *StudentService) CreateStudent(q *student.Student) error {
	q.Password, _ = s.generateUniquePassword(12)
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

	classDetail, err := s.classService.GetClassByID(q.ClassID)
	if err != nil {
		return nil, err
	}

	r = &student.RStudent{
		ID:              q.ID,
		Name:            q.Name,
		Phone:           q.Phone,
		Password:        q.Password,
		InstitutionID:   q.InstitutionID,
		InstitutionName: insDetail.Name,
		ClassID:         q.ClassID,
		ClassName:       classDetail.ClassName,
		Referrer:        q.Referrer,
		JobCode:         q.JobCode,
		JobName:         jobDetail.Name,
		JobDesc:         jobDetail.Desc,
		MajorIDs:        strings.Split(q.MajorIDs, ","),
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

		classDetail, err := s.classService.GetClassByID(v.ClassID)
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
			ClassName:       classDetail.ClassName,
			Referrer:        v.Referrer,
			JobCode:         v.JobCode,
			JobName:         jobDetail.Name,
			JobDesc:         jobDetail.Desc,
			MajorIDs:        strings.Split(v.MajorIDs, ","),
			MajorNames:      majorNames,
			Status:          v.Status,
			StatusName:      stable.RecordStatusMap[v.Status],
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
		})
	}

	return rStudents, total, nil
}

func (s *StudentService) GetStudentListBySortClass(page, pageSize int, keyword string, classId, institutionId, majorID int, status int) ([]student.RStudent, int64, error) {
	var (
		total1, total2 int64
		students       []student.Student
		rStudents      []student.RStudent
	)

	if classId <= 0 {
		return nil, 0, errors.New("class_id 不能为空")
	}

	db := global.GVA_DB.Model(&student.Student{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(job_code) LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	// 状态条件
	if status != 0 {
		db = db.Where("status =?", status)
	} else {
		db = db.Where("status != ?", stable.StatusDeleted)
	}

	// 分页计算
	offset := (page - 1) * pageSize

	// 获取 class_id 等于传递过来的 class_id 的数据集
	classDB := db.Where("class_id = ?", classId)
	if err := classDB.Count(&total1).Error; err != nil {
		return nil, 0, err
	}

	var classStudents []student.Student
	if err := classDB.Order("id DESC").Offset(offset).Limit(pageSize).Find(&classStudents).Error; err != nil {
		return nil, 0, err
	}

	// 计算剩余条目数
	remainingCount := pageSize - len(classStudents)
	if remainingCount > 0 {
		// 获取 major_id 为空的数据集
		nonClassDB := global.GVA_DB.Model(&student.Student{})
		if keyword != "" {
			searchQuery := "%" + strings.ToLower(keyword) + "%"
			nonClassDB = nonClassDB.Where("LOWER(name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(job_code) LIKE ?", searchQuery, searchQuery, searchQuery)
		}
		// 状态条件
		if status != 0 {
			nonClassDB = nonClassDB.Where("status =?", status)
		} else {
			nonClassDB = nonClassDB.Where("status != ?", stable.StatusDeleted)
		}
		nonClassDB = nonClassDB.Where("class_id != ?", classId)

		// 获取总数
		if err := nonClassDB.Count(&total2).Error; err != nil {
			return nil, 0, err
		}

		// 计算非 major 部分需要的偏移量
		nonClassOffset := max(0, int64(offset)-total1)
		nonClassStudents := []student.Student{}
		if err := nonClassDB.Order("id DESC").Offset(int(nonClassOffset)).Limit(remainingCount).Find(&nonClassStudents).Error; err != nil {
			return nil, 0, err
		}

		// 合并两个结果集
		students = append(classStudents, nonClassStudents...)
	} else {
		students = classStudents
	}

	for _, q := range students {
		// 获取 JobName 和 JobDesc
		jobDetail, err := s.jobService.GetJobByCode(q.JobCode)
		if err != nil {
			return nil, 0, err
		}
		insDetail, err := s.institutionService.GetInstitutionByID(q.InstitutionID)
		if err != nil {
			return nil, 0, err
		}
		classDetail, err := s.classService.GetClassByID(q.ClassID)
		if err != nil {
			return nil, 0, err
		}

		var classSorted int
		if q.ClassID == classId {
			classSorted = 1
		} else {
			classSorted = 0
		}
		rStudents = append(rStudents, student.RStudent{
			ID:              q.ID,
			Name:            q.Name,
			Phone:           q.Phone,
			Password:        q.Password,
			InstitutionID:   q.InstitutionID,
			InstitutionName: insDetail.Name,
			ClassID:         q.ClassID,
			ClassName:       classDetail.ClassName,
			Referrer:        q.Referrer,
			JobCode:         q.JobCode,
			JobName:         jobDetail.Name,
			JobDesc:         jobDetail.Desc,
			Status:          q.Status,
			StatusName:      stable.RecordStatusMap[q.Status],
			ClassSorted:     classSorted,
			CreateTime:      q.CreateTime,
			UpdateTime:      q.UpdateTime,
		})
	}

	return rStudents, total1 + total2, nil
}

func (s *StudentService) UpdateStudent(id int, q *student.Student) error {
	return s.DB.Model(&student.Student{}).Where("id = ?", id).Updates(q).Error
}

func (s *StudentService) BatchUpdateClass(studentIDs []int, classId int) error {
	if len(studentIDs) == 0 {
		return errors.New("student_ids 不能为空")
	}

	if classId == 0 {
		return errors.New("class_id 不能为0")
	}

	// 使用 IN 子句进行批量更新
	return s.DB.Model(&student.Student{}).
		Where("id IN (?)", studentIDs).
		Update("class_id", classId).Error
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

func (s *StudentService) generateUniquePassword(length int) (string, error) {
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
