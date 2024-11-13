package major

import (
	"admin_base_server/api/v1/major"
	qservice "admin_base_server/service/major"
	"fmt"
	"github.com/gin-gonic/gin"
)

// MajorRouter 路由管理器
type MajorRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/major [post,get,put,delete]
// @Router /admin/major-group [post,get,put,delete]
func (e *MajorRouter) InitMajorRouter(Router *gin.RouterGroup) {

	majorAPI = major.NewMajorAPI(qservice.NewMajorService())

	fmt.Println("majorAPI:", majorAPI)
	majorRouter := Router.Group("admin/major")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Major true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /major [post]
		majorRouter.POST("/major", majorAPI.CreateMajor)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Major
		// @Failure 400 {object} models.Response
		// @Router /major/{id} [get]
		majorRouter.GET("/major/:id", majorAPI.GetMajorByID)

		// @Summary 获取所有试题列表
		// @Description 获取所有试题列表
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Major
		// @Failure 400 {object} models.Response
		// @Router /major/list [get]
		majorRouter.GET("/major/list", majorAPI.GetMajorList)

		// @Summary 更新试题
		// @Description 根据ID更新试题信息
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Param body body models.Major true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /major/{id} [put]
		majorRouter.PUT("/major/:id", majorAPI.UpdateMajor)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /major/{id} [delete]
		majorRouter.DELETE("/major/:id", majorAPI.DeleteMajor)

		// @Summary 批量导入试题
		// @Description 批量导入试题
		// @Tags 试题管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /major/batch-import [post]
		majorRouter.POST("/major/batch-import", majorAPI.BatchImportMajors)
	}
}
