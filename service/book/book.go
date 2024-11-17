package book

import (
	"admin_base_server/global"
	"admin_base_server/model/book"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	"admin_base_server/service/topic"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type BookService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
	topicService  *topic.TopicService
}

func NewBookService() *BookService {
	return &BookService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
		topicService:  topic.NewTopicService(),
	}
}

func (s *BookService) CreateBook(q *book.RBook) error {
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

	questions, questionsNumber, err := s.generateBook(q.Level, q.MajorID, q.Component, q.UnitNumber)
	if err != nil {
		return err
	}

	bookModel := book.Book{
		ID:              q.ID,
		Name:            q.Name,
		MajorID:         q.MajorID,
		Level:           q.Level,
		Component:       q.Component,
		UnitNumber:      q.UnitNumber,
		Questions:       questions,
		QuestionsNumber: questionsNumber,
		Creator:         q.Creator,
		TemplateID:      q.TemplateID,
		Tag:             q.Tag,
	}

	return s.DB.Create(bookModel).Error
}

func (s *BookService) GetBookByID(id int) (*book.RBook, error) {
	var (
		q *book.Book
		r *book.RBook
	)
	if err := s.DB.First(&q, id).Error; err != nil {
		return nil, err
	}
	_, levelMap, err := s.getConfigMap()
	if err != nil {
		return nil, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, err
	}

	var (
		component []*book.Component
		questions [][]*book.Questions
	)
	err = json.Unmarshal([]byte(cast.ToString(q.Component)), &component)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cast.ToString(q.Questions)), &questions)
	if err != nil {
		return nil, err
	}

	r = &book.RBook{
		ID:              q.ID,
		Name:            q.Name,
		MajorID:         q.MajorID,
		MajorName:       majorMap[q.MajorID],
		Level:           q.Level,
		LevelName:       levelMap[q.Level],
		Component:       component,
		UnitNumber:      q.UnitNumber,
		Questions:       questions,
		QuestionsNumber: q.QuestionsNumber,
		Creator:         q.Creator,
		TemplateID:      q.TemplateID,
		TemplateName:    "demo",
		Tag:             q.Tag,
		CreateTime:      q.CreateTime,
		UpdateTime:      q.UpdateTime,
	}
	return r, nil
}

func (s *BookService) GetBookList(page, pageSize int, keyword, level string, majorID int) ([]book.RBook, int64, error) {
	var (
		total  int64
		books  []book.Book
		rBooks []book.RBook
	)

	db := global.GVA_DB.Model(&book.Book{})

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
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	_, levelMap, err := s.getConfigMap()
	if err != nil {
		return nil, 0, err
	}

	majorMap, err := s.getMajorMap()
	if err != nil {
		return nil, 0, err
	}

	for _, v := range books {
		var (
			component []*book.Component
			questions [][]*book.Questions
		)
		err = json.Unmarshal([]byte(cast.ToString(v.Component)), &component)
		if err != nil {
			return nil, 0, err
		}
		err = json.Unmarshal([]byte(cast.ToString(v.Questions)), &questions)
		if err != nil {
			return nil, 0, err
		}

		rBooks = append(rBooks, book.RBook{
			ID:              v.ID,
			Name:            v.Name,
			MajorID:         v.MajorID,
			MajorName:       majorMap[v.MajorID],
			Level:           v.Level,
			LevelName:       levelMap[v.Level],
			Component:       component,
			UnitNumber:      v.UnitNumber,
			Questions:       questions,
			QuestionsNumber: v.QuestionsNumber,
			Creator:         v.Creator,
			TemplateID:      v.TemplateID,
			TemplateName:    "demo",
			Tag:             v.Tag,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
		})
	}

	return rBooks, total, nil
}

func (s *BookService) DeleteBook(id int) error {
	return s.DB.Delete(&book.Book{}, id).Error
}

func (s *BookService) getConfigMap() (map[string]string, map[string]string, error) {
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

func (s *BookService) getMajorMap() (map[int]string, error) {
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

func (s *BookService) generateBook(level string, majorId int, component []*book.Component, number int) (questions [][]*book.Questions, questionNumber int, err error) {
	var (
		questionIds []int
	)
	for i := 0; i < number; i++ {
		var page []*book.Questions
		for _, v := range component {
			var ids []int
			qs, _, err := s.topicService.GetTopicList(1, v.Number, "", v.Key, level, majorId, questionIds)
			if err != nil {
				return questions, questionNumber, err
			}
			for _, q := range qs {
				ids = append(ids, q.ID)
			}
			page = append(page, &book.Questions{
				Key: v.Key,
				Ids: ids,
			})
			questionIds = append(questionIds, ids...)
		}
		questions = append(questions, page)
	}
	questionNumber = len(questionIds)

	return
}
