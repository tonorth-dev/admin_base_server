package job

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/job"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type JobService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
}

func NewJobService() *JobService {
	return &JobService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
	}
}

func (s *JobService) CreateJob(q *job.Job) error {
	q.Status = stable.StatusActive

	return s.DB.Create(q).Error
}

func (s *JobService) GetJobByID(id int) (*job.RJob, error) {
	var (
		q *job.Job
		r *job.RJob
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, err
	}

	tmpMap := map[string]string{}
	_ = json.Unmarshal(q.Condition, &tmpMap)

	r = &job.RJob{
		ID:              q.ID,
		Code:            q.Code,
		Name:            q.Name,
		Desc:            q.Desc,
		Cate:            q.Cate,
		CompanyCode:     q.CompanyCode,
		CompanyName:     q.CompanyName,
		EnrollmentNum:   q.EnrollmentNum,
		EnrollmentRatio: q.EnrollmentRatio,
		Condition:       q.Condition,
		ConditionName:   generateConditionDesc(tmpMap),
		MajorID:         q.MajorID,
		MajorName:       majorMap[q.MajorID],
		City:            q.City,
		Phone:           q.Phone,
		Status:          q.Status,
		StatusName:      stable.RecordStatusMap[q.Status],
		CreateTime:      q.CreateTime,
		UpdateTime:      q.UpdateTime,
	}
	return r, nil
}

func (s *JobService) GetJobList(page, pageSize int, keyword string, majorID int) ([]job.RJob, int64, error) {
	var (
		total int64
		jobs  []job.Job
		rJobs []job.RJob
	)

	db := global.GVA_DB.Model(&job.Job{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(company_code) LIKE ? OR LOWER(company_name) LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery)
	}

	// 筛选条件
	if majorID != 0 {
		db = db.Where("major_id = ?", majorID)
	}
	db = db.Where("status = ?", stable.StatusActive)

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, 0, err
	}

	for _, q := range jobs {
		tmpMap := map[string]string{}
		_ = json.Unmarshal(q.Condition, &tmpMap)
		rJobs = append(rJobs, job.RJob{
			ID:              q.ID,
			Code:            q.Code,
			Name:            q.Name,
			Desc:            q.Desc,
			Cate:            q.Cate,
			CompanyCode:     q.CompanyCode,
			CompanyName:     q.CompanyName,
			EnrollmentNum:   q.EnrollmentNum,
			EnrollmentRatio: q.EnrollmentRatio,
			Condition:       q.Condition,
			ConditionName:   generateConditionDesc(tmpMap),
			MajorID:         q.MajorID,
			MajorName:       majorMap[q.MajorID],
			City:            q.City,
			Phone:           q.Phone,
			Status:          q.Status,
			StatusName:      stable.RecordStatusMap[q.Status],
			CreateTime:      q.CreateTime,
			UpdateTime:      q.UpdateTime,
		})
	}

	return rJobs, total, nil
}

func (s *JobService) UpdateJob(id int, q *job.Job) error {
	return s.DB.Model(&job.Job{}).Where("id = ?", id).Updates(q).Error
}

func (s *JobService) DeleteJob(ids []int) error {
	return s.DB.Model(&job.Job{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
}

// todo 加上校验逻辑
func (s *JobService) BatchImportJobs(jobs []job.Job) error {
	return s.DB.Create(&jobs).Error
}

func (s *JobService) ExportJobs() ([]job.Job, error) {
	var jobs []job.Job
	if err := s.DB.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *JobService) getConfigMap() (map[string]string, map[string]string, error) {
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

func (s *JobService) getMajorMap() (map[int]string, error) {
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

// 根据ID查找对应的名称
func findNameById(data []map[string]string, id string) string {
	for _, item := range data {
		if item["id"] == id {
			return item["name"]
		}
	}
	return id // 如果找不到，返回原始ID
}

// 生成condition_desc
func generateConditionDesc(condition map[string]string) string {
	var desc string

	// 添加source
	sourceName := findNameById(job.Source, condition["source"])
	desc += fmt.Sprintf("来源: %s\t", sourceName)

	// 添加degree
	degreeName := findNameById(job.Degree, condition["degree"])
	desc += fmt.Sprintf("学位: %s\t", degreeName)

	// 添加qualification
	qualificationName := findNameById(job.Qualification, condition["qualification"])
	desc += fmt.Sprintf("学历: %s\n", qualificationName)

	// 添加exam
	desc += fmt.Sprintf("考试专业科目: %s\n", condition["exam"])

	// 添加major
	desc += fmt.Sprintf("所学专业: %s\n", condition["major"])

	// 添加other
	desc += fmt.Sprintf("其他条件: %s", condition["other"])

	return desc
}