package institution

import (
	"admin_base_server/api/v1/institution"
	qservice "admin_base_server/service/institution"
	"github.com/gin-gonic/gin"
)

// InstitutionRouter 路由管理器
type InstitutionRouter struct{}

// @Summary 初始化机构和题组路由
// @Description 初始化机构和题组的相关路由
// @Tags admin
// @Router /admin/institution [post,get,put,delete]
// @Router /admin/institution-group [post,get,put,delete]
func (e *InstitutionRouter) InitInstitutionRouter(Router *gin.RouterGroup) {

	institutionAPI = institution.NewInstitutionAPI(qservice.NewInstitutionService())

	institutionRouter := Router.Group("admin/institution")
	{
		// 机构路由
		// @Summary 创建机构
		// @Description 创建一个新的机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param body body models.Institution true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /institution [post]
		institutionRouter.POST("/institution", institutionAPI.CreateInstitution)

		// @Summary 获取机构详情
		// @Description 根据ID获取机构详情
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Institution
		// @Failure 400 {object} models.Response
		// @Router /institution/{id} [get]
		institutionRouter.GET("/institution/:id", institutionAPI.GetInstitutionByID)

		// @Summary 获取所有机构列表
		// @Description 获取所有机构列表
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Institution
		// @Failure 400 {object} models.Response
		// @Router /institution/list [get]
		institutionRouter.GET("/institution/list", institutionAPI.GetInstitutionList)

		// @Summary 更新机构
		// @Description 根据ID更新机构信息
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Param body body models.Institution true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /institution/{id} [put]
		institutionRouter.PUT("/institution/:id", institutionAPI.UpdateInstitution)

		// @Summary 删除机构
		// @Description 根据ID删除机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /institution/{id} [delete]
		institutionRouter.DELETE("/institution/:id", institutionAPI.DeleteInstitution)

		// @Summary 批量导入机构
		// @Description 批量导入机构
		// @Tags 机构管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /institution/batch-import [post]
		institutionRouter.POST("/institution/batch-import", institutionAPI.BatchImportInstitutions)

		// @Summary 导出机构
		// @Description 导出机构
		// @Tags 机构管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /institution/export [get]
		institutionRouter.GET("/institution/export", institutionAPI.ExportInstitutions)
	}
}
