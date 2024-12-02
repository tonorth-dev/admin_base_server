package job

import (
	"admin_base_server/api/v1/job"
	qservice "admin_base_server/service/job"
	"github.com/gin-gonic/gin"
)

// JobRouter 路由管理器
type JobRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/job [post,get,put,delete]
// @Router /admin/job-group [post,get,put,delete]
func (e *JobRouter) InitJobRouter(Router *gin.RouterGroup) {

	jobAPI = job.NewJobAPI(qservice.NewJobService())

	jobRouter := Router.Group("admin/job")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Job true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /job [post]
		jobRouter.POST("/job", jobAPI.CreateJob)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Job
		// @Failure 400 {object} models.Response
		// @Router /job/{id} [get]
		jobRouter.GET("/job/:id", jobAPI.GetJobByID)

		// @Summary 获取所有试题列表
		// @Description 获取所有试题列表
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Job
		// @Failure 400 {object} models.Response
		// @Router /job/list [get]
		jobRouter.GET("/job/list", jobAPI.GetJobList)

		// @Summary 更新试题
		// @Description 根据ID更新试题信息
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Param body body models.Job true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /job/{id} [put]
		jobRouter.PUT("/job/:id", jobAPI.UpdateJob)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /job/{id} [delete]
		jobRouter.DELETE("/job/:id", jobAPI.DeleteJob)

		// @Summary 批量导入试题
		// @Description 批量导入试题
		// @Tags 试题管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /job/batch-import [post]
		jobRouter.POST("/job/batch-import", jobAPI.BatchImportJobs)

		// @Summary 导出试题
		// @Description 导出试题
		// @Tags 试题管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /job/export [get]
		jobRouter.GET("/job/export", jobAPI.ExportJobs)
	}
}
