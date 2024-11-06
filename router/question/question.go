package question

import (
	"github.com/gin-gonic/gin"
)

// QuestionRouter 路由管理器
type QuestionRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/question [post,get,put,delete]
// @Router /admin/question-group [post,get,put,delete]
func (e *QuestionRouter) InitQuestionRouter(Router *gin.RouterGroup) {
	questionRouter := Router.Group("admin/question")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Question true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question [post]
		questionRouter.POST("/question", questionAPI.CreateQuestion)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Question
		// @Failure 400 {object} models.Response
		// @Router /question/{id} [get]
		questionRouter.GET("/question/:id", questionAPI.GetQuestionByID)

		// @Summary 更新试题
		// @Description 根据ID更新试题信息
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Param body body models.Question true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question/{id} [put]
		questionRouter.PUT("/question/:id", questionAPI.UpdateQuestion)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question/{id} [delete]
		questionRouter.DELETE("/question/:id", questionAPI.DeleteQuestion)

		// @Summary 批量导入试题
		// @Description 批量导入试题
		// @Tags 试题管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question/batch-import [post]
		questionRouter.POST("/question/batch-import", questionAPI.BatchImportQuestions)

		// @Summary 导出试题
		// @Description 导出试题
		// @Tags 试题管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /question/export [get]
		questionRouter.GET("/question/export", questionAPI.ExportQuestions)
	}
	{
		// 题组路由
		// @Summary 创建题组
		// @Description 创建一个新的题组
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param body body models.QuestionGroup true "题组信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question-group [post]
		questionRouter.POST("/question-group", questionGroupAPI.CreateQuestionGroup)

		// @Summary 获取题组详情
		// @Description 根据ID获取题组详情
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Success 200 {object} models.QuestionGroup
		// @Failure 400 {object} models.Response
		// @Router /question-group/{id} [get]
		questionRouter.GET("/question-group/:id", questionGroupAPI.GetQuestionGroupByID)

		// @Summary 更新题组
		// @Description 根据ID更新题组信息
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Param body body models.QuestionGroup true "题组信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question-group/{id} [put]
		questionRouter.PUT("/question-group/:id", questionGroupAPI.UpdateQuestionGroup)

		// @Summary 删除题组
		// @Description 根据ID删除题组
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /question-group/{id} [delete]
		questionRouter.DELETE("/question-group/:id", questionGroupAPI.DeleteQuestionGroup)

		// @Summary 导出题组为PDF
		// @Description 根据ID导出题组为PDF
		// @Tags 题组管理
		// @Accept json
		// @Produce octet-stream
		// @Param id path string true "题组ID"
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /question-group/export-pdf/{id} [get]
		questionRouter.GET("/question-group/export-pdf/:id", questionGroupAPI.ExportQuestionGroupPDF)
	}
}
