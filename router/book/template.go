package book

import (
	"admin_base_server/api/v1/book"
	qservice "admin_base_server/service/book"
	"github.com/gin-gonic/gin"
)

// TemplateRouter 路由管理器
type TemplateRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/template [post,get,put,delete]
// @Router /admin/template-group [post,get,put,delete]
func (e *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup) {

	templateAPI = book.NewTemplateAPI(qservice.NewTemplateService())

	templateRouter := Router.Group("admin/book")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Template true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /template [post]
		templateRouter.POST("/template", templateAPI.CreateTemplate)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Template
		// @Failure 400 {object} models.Response
		// @Router /template/{id} [get]
		templateRouter.GET("/template/:id", templateAPI.GetTemplateByID)

		// @Summary 获取所有试题列表
		// @Description 获取所有试题列表
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Template
		// @Failure 400 {object} models.Response
		// @Router /template/list [get]
		templateRouter.GET("/template/list", templateAPI.GetTemplateList)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /template/{id} [delete]
		templateRouter.DELETE("/template/:id", templateAPI.DeleteTemplate)
	}
}
