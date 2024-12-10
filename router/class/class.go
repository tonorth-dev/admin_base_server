package class

import (
	"admin_base_server/api/v1/class"
	qservice "admin_base_server/service/class"
	"github.com/gin-gonic/gin"
)

// ClassRouter 路由管理器
type ClassRouter struct{}

// @Summary 初始化机构和题组路由
// @Description 初始化机构和题组的相关路由
// @Tags admin
// @Router /admin/class [post,get,put,delete]
// @Router /admin/class-group [post,get,put,delete]
func (e *ClassRouter) InitClassRouter(Router *gin.RouterGroup) {

	classAPI = class.NewClassAPI(qservice.NewClassService())

	classRouter := Router.Group("admin/class")
	{
		// 机构路由
		// @Summary 创建机构
		// @Description 创建一个新的机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param body body models.Class true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /class [post]
		classRouter.POST("/class", classAPI.CreateClass)

		// @Summary 获取机构详情
		// @Description 根据ID获取机构详情
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Class
		// @Failure 400 {object} models.Response
		// @Router /class/{id} [get]
		classRouter.GET("/class/:id", classAPI.GetClassByID)

		// @Summary 获取所有机构列表
		// @Description 获取所有机构列表
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Class
		// @Failure 400 {object} models.Response
		// @Router /class/list [get]
		classRouter.GET("/class/list", classAPI.GetClassList)

		// @Summary 更新机构
		// @Description 根据ID更新机构信息
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Param body body models.Class true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /class/{id} [put]
		classRouter.PUT("/class/:id", classAPI.UpdateClass)

		// @Summary 删除机构
		// @Description 根据ID删除机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /class/{id} [delete]
		classRouter.DELETE("/class/:id", classAPI.DeleteClass)

		// @Summary 批量导入机构
		// @Description 批量导入机构
		// @Tags 机构管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /class/batch-import [post]
		classRouter.POST("/class/batch-import", classAPI.BatchImportClasss)

		// @Summary 导出机构
		// @Description 导出机构
		// @Tags 机构管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /class/export [get]
		classRouter.GET("/class/export", classAPI.ExportClasss)
	}
}
