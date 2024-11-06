package question

import (
	"admin_base_server/api/v1/question"
	qservice "admin_base_server/service/question"
)

type RouterGroup struct {
	QuestionRouter
}

var (
	questionAPI      = question.NewQuestionAPI(qservice.NewQuestionService())
	questionGroupAPI = question.NewQuestionGroupAPI(qservice.NewQuestionGroupService())
)
