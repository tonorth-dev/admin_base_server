package topic

import (
	"admin_base_server/global"
	stable "admin_base_server/model/const"
	"admin_base_server/model/topic"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type TopicService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
}

func NewTopicService() *TopicService {
	return &TopicService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
	}
}

func (s *TopicService) CreateTopic(q *topic.Topic) error {

	cateMap, levelMap, err := s.getConfigMap()
	if err != nil {
		return err
	}

	if _, found := levelMap[q.Level]; !found {
		return errors.New("问题等级未定义")
	}

	if _, found := cateMap[q.Cate]; !found {
		return errors.New("问题类型未定义")
	}

	return s.DB.Create(q).Error
}

func (s *TopicService) GetTopicByID(id int) (*topic.RTopic, error) {
	var (
		q *topic.Topic
		r *topic.RTopic
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

	r = &topic.RTopic{
		ID:          q.ID,
		Title:       q.Title,
		Author:      q.Author,
		Answer:      q.Answer,
		AnswerDraft: q.AnswerDraft,
		Tag:         q.Tag,
		Level:       q.Level,
		LevelName:   levelMap[q.Level],
		Cate:        q.Cate,
		CateName:    cateMap[q.Cate],
		MajorID:     q.MajorID,
		MajorName:   majorMap[q.MajorID],
		Status:      q.Status,
		StatusName:  stable.RecordStatusMap[q.Status],
		CreateTime:  q.CreateTime,
		UpdateTime:  q.UpdateTime,
	}
	return r, nil
}

func (s *TopicService) GetTopicList(page, pageSize int, keyword, cate, level string, majorID int, status int) ([]topic.RTopic, int64, error) {
	var (
		total   int64
		topics  []topic.Topic
		rTopics []topic.RTopic
	)

	db := global.GVA_DB.Model(&topic.Topic{})

	// 搜索条件
	if keyword != "" {
		searchQuery := "%" + strings.ToLower(keyword) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(answer) LIKE ?", searchQuery, searchQuery, searchQuery)
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
	if status == 0 {
		db = db.Where("status in (?)", []int{stable.StatusActive, stable.StatusDraft, stable.StatusAuditing})
	} else {
		db = db.Where("status = ?", status)
	}

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&topics).Error; err != nil {
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

	for _, v := range topics {
		rTopics = append(rTopics, topic.RTopic{
			ID:         v.ID,
			Title:      v.Title,
			Author:     v.Author,
			Answer:     v.Answer,
			Tag:        v.Tag,
			Level:      v.Level,
			LevelName:  levelMap[v.Level],
			Cate:       v.Cate,
			CateName:   cateMap[v.Cate],
			MajorID:    v.MajorID,
			MajorName:  majorMap[v.MajorID],
			Status:     v.Status,
			StatusName: stable.RecordStatusMap[v.Status],
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		})
	}

	return rTopics, total, nil
}

func (s *TopicService) UpdateTopic(id int, q *topic.Topic) error {
	return s.DB.Model(&topic.Topic{}).Where("id = ?", id).Updates(q).Error
}

func (s *TopicService) SubmitAudit(id int, a *topic.Audit) error {
	r, err := s.GetTopicByID(id)
	if err != nil {
		return errors.New("试题不存在")
	}

	if r.Status != stable.StatusDraft {
		return errors.New("试题不在草稿状态")
	}

	q := &topic.Topic{
		Answer:  a.Answer,
		Invitee: a.Invitee,
		Status:  stable.StatusAuditing,
	}
	return s.UpdateTopic(id, q)
}

func (s *TopicService) AuditTopic(id int, approved bool) error {
	r, err := s.GetTopicByID(id)
	if err != nil {
		return errors.New("试题不存在")
	}

	if r.Status != stable.StatusAuditing {
		return errors.New("试题不在审核中状态")
	}

	q := &topic.Topic{}
	if approved {
		q.Status = stable.StatusActive
	} else {
		q.Status = stable.StatusDraft
		q.Answer = r.AnswerDraft
	}

	return s.UpdateTopic(id, q)
}

func (s *TopicService) DeleteTopic(ids []int) error {
	return s.DB.Model(&topic.Topic{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
}

// todo 加上校验逻辑
func (s *TopicService) BatchImportTopics(topics []topic.Topic) error {
	return s.DB.Create(&topics).Error
}

func (s *TopicService) ExportTopics() ([]topic.Topic, error) {
	var topics []topic.Topic
	if err := s.DB.Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

func (s *TopicService) getConfigMap() (map[string]string, map[string]string, error) {
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

func (s *TopicService) getMajorMap() (map[int]string, error) {
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
