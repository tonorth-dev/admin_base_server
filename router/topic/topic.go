package topic

import (
	"admin_base_server/api/v1/topic"
	qservice "admin_base_server/service/topic"
	"github.com/gin-gonic/gin"
)

// TopicRouter 路由管理器
type TopicRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/topic [post,get,put,delete]
// @Router /admin/topic-group [post,get,put,delete]
func (e *TopicRouter) InitTopicRouter(Router *gin.RouterGroup) {

	topicAPI = topic.NewTopicAPI(qservice.NewTopicService())

	topicRouter := Router.Group("admin/topic")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Topic true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /topic [post]
		topicRouter.POST("/topic", topicAPI.CreateTopic)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Topic
		// @Failure 400 {object} models.Response
		// @Router /topic/{id} [get]
		topicRouter.GET("/topic/:id", topicAPI.GetTopicByID)

		// @Summary 获取所有试题列表
		// @Description 获取所有试题列表
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Topic
		// @Failure 400 {object} models.Response
		// @Router /topic/list [get]
		topicRouter.GET("/topic/list", topicAPI.GetTopicList)

		// @Summary 更新试题
		// @Description 根据ID更新试题信息
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Param body body models.Topic true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /topic/{id} [put]
		topicRouter.PUT("/topic/:id", topicAPI.UpdateTopic)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /topic/{id} [delete]
		topicRouter.DELETE("/topic/:id", topicAPI.DeleteTopic)

		// @Summary 批量导入试题
		// @Description 批量导入试题
		// @Tags 试题管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /topic/batch-import [post]
		topicRouter.POST("/topic/batch-import", topicAPI.BatchImportTopics)

		// @Summary 导出试题
		// @Description 导出试题
		// @Tags 试题管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /topic/export [get]
		topicRouter.GET("/topic/export", topicAPI.ExportTopics)
	}
}
