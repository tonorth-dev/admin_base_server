package student

import (
	"admin_base_server/api/v1/student"
	qservice "admin_base_server/service/student"
	"github.com/gin-gonic/gin"
)

// StudentRouter 路由管理器
type StudentRouter struct{}

// @Summary 初始化机构和题组路由
// @Description 初始化机构和题组的相关路由
// @Tags admin
// @Router /admin/student [post,get,put,delete]
// @Router /admin/student-group [post,get,put,delete]
func (e *StudentRouter) InitStudentRouter(Router *gin.RouterGroup) {

	studentAPI = student.NewStudentAPI(qservice.NewStudentService())

	studentRouter := Router.Group("admin/student")
	{
		// 机构路由
		// @Summary 创建机构
		// @Description 创建一个新的机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param body body models.Student true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /student [post]
		studentRouter.POST("/student", studentAPI.CreateStudent)

		// @Summary 获取机构详情
		// @Description 根据ID获取机构详情
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Student
		// @Failure 400 {object} models.Response
		// @Router /student/{id} [get]
		studentRouter.GET("/student/:id", studentAPI.GetStudentByID)

		// @Summary 获取所有机构列表
		// @Description 获取所有机构列表
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Student
		// @Failure 400 {object} models.Response
		// @Router /student/list [get]
		studentRouter.GET("/student/list", studentAPI.GetStudentList)

		// @Summary 更新机构
		// @Description 根据ID更新机构信息
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Param body body models.Student true "机构信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /student/{id} [put]
		studentRouter.PUT("/student/:id", studentAPI.UpdateStudent)

		// @Summary 删除机构
		// @Description 根据ID删除机构
		// @Tags 机构管理
		// @Accept json
		// @Produce json
		// @Param id path string true "机构ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /student/{id} [delete]
		studentRouter.DELETE("/student/:id", studentAPI.DeleteStudent)

		// @Summary 批量导入机构
		// @Description 批量导入机构
		// @Tags 机构管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /student/batch-import [post]
		studentRouter.POST("/student/batch-import", studentAPI.BatchImportStudents)

		// @Summary 导出机构
		// @Description 导出机构
		// @Tags 机构管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /student/export [get]
		studentRouter.GET("/student/export", studentAPI.ExportStudents)
	}
}
