package book

import (
	"admin_base_server/global"
	"admin_base_server/model/book"
	stable "admin_base_server/model/const"
	"admin_base_server/model/topic"
	"admin_base_server/service/config"
	"admin_base_server/service/major"
	stopic "admin_base_server/service/topic"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

type BookService struct {
	DB            *gorm.DB
	configService *config.ConfigService
	majorService  *major.MajorService
	topicService  *stopic.TopicService
}

func NewBookService() *BookService {
	return &BookService{
		DB:            global.GVA_DB,
		configService: config.NewConfigService(),
		majorService:  major.NewMajorService(),
		topicService:  stopic.NewTopicService(),
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

	componentJSON, _ := json.Marshal(q.Component)
	questionsJSON, _ := json.Marshal(questions)

	bookModel := &book.Book{
		ID:              q.ID,
		Name:            q.Name,
		MajorID:         q.MajorID,
		Level:           q.Level,
		Component:       componentJSON,
		UnitNumber:      q.UnitNumber,
		Questions:       questionsJSON,
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
	cateMap, levelMap, err := s.getConfigMap()
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
	componentDesc, _ := s.generateComponentDesc(component, cateMap)
	questionsDesc, _ := s.generateQuestionsDesc(questions, cateMap)

	r = &book.RBook{
		ID:              q.ID,
		Name:            q.Name,
		MajorID:         q.MajorID,
		MajorName:       majorMap[q.MajorID],
		Level:           q.Level,
		LevelName:       levelMap[q.Level],
		Component:       component,
		ComponentDesc:   componentDesc,
		UnitNumber:      q.UnitNumber,
		Questions:       questions,
		QuestionsNumber: q.QuestionsNumber,
		QuestionsDesc:   questionsDesc,
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
	if err := db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&books).Error; err != nil {
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
		componentDesc, _ := s.generateComponentDesc(component, cateMap)
		//questionsDesc, _ := s.generateQuestionsDesc(questions, cateMap)

		rBooks = append(rBooks, book.RBook{
			ID:            v.ID,
			Name:          v.Name,
			MajorID:       v.MajorID,
			MajorName:     majorMap[v.MajorID],
			Level:         v.Level,
			LevelName:     levelMap[v.Level],
			Component:     component,
			ComponentDesc: componentDesc,
			UnitNumber:    v.UnitNumber,
			Questions:     questions,
			//QuestionsDesc:   questionsDesc,
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

func (s *BookService) DeleteBook(ids []int) error {
	return s.DB.Model(&book.Book{}).Where("id IN (?)", ids).Update("status", stable.StatusDeleted).Error
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

func (s *BookService) generateBook(level string, majorId int, component []*book.Component, number int) (questions [][]*book.Questions, questionNumber int, err error) {
	var (
		questionIds []int
		wg          sync.WaitGroup
		mu          sync.Mutex
		errs        []error
	)

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var page []*book.Questions
			localQuestionIds := make([]int, len(questionIds))
			copy(localQuestionIds, questionIds)

			for _, v := range component {
				var ids []int
				qs, _, err := s.topicService.GetTopicList(1, v.Number, "", v.Key, level, majorId, stable.StatusActive)
				if err != nil {
					mu.Lock()
					errs = append(errs, err)
					mu.Unlock()
					return
				}
				for _, q := range qs {
					ids = append(ids, q.ID)
				}
				page = append(page, &book.Questions{
					Key: v.Key,
					Ids: ids,
				})
				localQuestionIds = append(localQuestionIds, ids...)
			}

			mu.Lock()
			questionIds = append(questionIds, localQuestionIds[len(questionIds):]...)
			questions = append(questions, page)
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	if len(errs) > 0 {
		return questions, questionNumber, errors.New("one or more errors occurred during processing")
	}

	questionNumber = len(questionIds)
	return
}

func (s *BookService) generateComponentDesc(component []*book.Component, catMap map[string]string) ([]string, error) {
	var (
		desc []string
	)
	for _, v := range component {
		desc = append(desc, catMap[v.Key]+"：数量"+strconv.Itoa(v.Number))
	}

	return desc, nil
}

func (s *BookService) generateQuestionsDesc(questions [][]*book.Questions, catMap map[string]string) (qds []*book.QuestionsDetails, err error) {
	if len(questions) == 0 {
		return nil, errors.New("练习题目为空")
	}
	for key, val := range questions {
		var qs []*book.QuestionsDetail
		for _, v := range val {
			qd := &book.QuestionsDetail{
				CateName: catMap[v.Key],
				List:     make([]*topic.RTopic, 0),
			}
			for _, id := range v.Ids {
				q, err := s.topicService.GetTopicByID(id)
				if err != nil {
					return nil, err
				}
				qd.List = append(qd.List, q)
			}
			qs = append(qs, qd)
		}
		qds = append(qds, &book.QuestionsDetails{
			Title: fmt.Sprintf("第%d节", key+1),
			Data:  qs,
		})
	}

	return
}
