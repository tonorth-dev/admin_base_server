package topic

import (
	"admin_base_server/global"
	"admin_base_server/model/topic"
	"gorm.io/gorm"
	"strings"
)

type TopicService struct {
	DB *gorm.DB
}

func NewTopicService() *TopicService {
	return &TopicService{DB: global.GVA_DB}
}

func (s *TopicService) CreateTopic(q *topic.Topic) error {
	return s.DB.Create(q).Error
}

func (s *TopicService) GetTopicByID(id int) (*topic.Topic, error) {
	var q topic.Topic
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}
	return &q, nil
}

func (s *TopicService) GetTopicList(page, pageSize int, search string, cate, majorID int) ([]topic.Topic, int64, error) {
	var topics []topic.Topic
	var total int64

	db := global.GVA_DB.Model(&topic.Topic{})

	// 搜索条件
	if search != "" {
		searchQuery := "%" + strings.ToLower(search) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(author) LIKE ? OR LOWER(answer) LIKE ?", searchQuery, searchQuery, searchQuery)
	}

	// 筛选条件
	if cate != 0 {
		db = db.Where("cate = ?", cate)
	}
	if majorID != 0 {
		db = db.Where("major_id = ?", majorID)
	}

	// 分页
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, 0, err
	}

	return topics, total, nil
}

func (s *TopicService) UpdateTopic(id int, q *topic.Topic) error {
	return s.DB.Model(&topic.Topic{}).Where("id = ?", id).Updates(q).Error
}

func (s *TopicService) DeleteTopic(id int) error {
	return s.DB.Delete(&topic.Topic{}, id).Error
}

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
