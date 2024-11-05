package example

import (
	"admin_base_server/api/v1/question"
	qservice "admin_base_server/service/question"
	"github.com/gin-gonic/gin"
)

type QuestionRouter struct{}

func (e *QuestionRouter) InitQuestionRouter(Router *gin.RouterGroup) {
	questionRouter := Router.Group("question")
	questionAPI := question.QuestionAPI{
		Service: qservice.NewQuestionService(), // Replace with your actual service instance
	}
	questionGroupAPI := question.QuestionGroupAPI{
		Service: qservice.NewQuestionGroupService(), // Replace with your actual service instance
	}
	{
		// 试题路由
		questionRouter.POST("/question", questionAPI.CreateQuestion)
		questionRouter.GET("/question/:id", questionAPI.GetQuestionByID)
		questionRouter.PUT("/question/:id", questionAPI.UpdateQuestion)
		questionRouter.DELETE("/question/:id", questionAPI.DeleteQuestion)
		questionRouter.POST("/question/batch-import", questionAPI.BatchImportQuestions)
		questionRouter.GET("/question/export", questionAPI.ExportQuestions)
	}
	{
		// 题组路由
		questionRouter.POST("/question-group", questionGroupAPI.CreateQuestionGroup)
		questionRouter.GET("/question-group/:id", questionGroupAPI.GetQuestionGroupByID)
		questionRouter.PUT("/question-group/:id", questionGroupAPI.UpdateQuestionGroup)
		questionRouter.DELETE("/question-group/:id", questionGroupAPI.DeleteQuestionGroup)
		questionRouter.GET("/question-group/export-pdf/:id", questionGroupAPI.ExportQuestionGroupPDF) // 删除切片
	}
}
