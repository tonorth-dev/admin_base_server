package question

import (
	"admin_base_server/global"
	"admin_base_server/model/question"
	"gorm.io/gorm"
)

type QuestionService struct {
	DB *gorm.DB
}

func NewQuestionService() *QuestionService {
	return &QuestionService{DB: global.GVA_DB}
}

func (s *QuestionService) CreateQuestion(q *question.Question) error {
	return s.DB.Create(q).Error
}

func (s *QuestionService) GetQuestionByID(id int) (*question.Question, error) {
	var q question.Question
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}
	return &q, nil
}

func (s *QuestionService) UpdateQuestion(id int, q *question.Question) error {
	return s.DB.Model(&question.Question{}).Where("id = ?", id).Updates(q).Error
}

func (s *QuestionService) DeleteQuestion(id int) error {
	return s.DB.Delete(&question.Question{}, id).Error
}

func (s *QuestionService) BatchImportQuestions(questions []question.Question) error {
	return s.DB.Create(&questions).Error
}

func (s *QuestionService) ExportQuestions() ([]question.Question, error) {
	var questions []question.Question
	if err := s.DB.Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}
