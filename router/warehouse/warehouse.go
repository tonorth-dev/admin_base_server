package warehouse

import (
	"admin_base_server/api/v1/warehouse"
	qservice "admin_base_server/service/warehouse"
	"fmt"
	"github.com/gin-gonic/gin"
)

// WarehouseRouter 路由管理器
type WarehouseRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/warehouse [post,get,put,delete]
// @Router /admin/warehouse-group [post,get,put,delete]
func (e *WarehouseRouter) InitWarehouseRouter(Router *gin.RouterGroup) {

	warehouseAPI = warehouse.NewWarehouseAPI(qservice.NewWarehouseService())
	warehouseGroupAPI = warehouse.NewWarehouseGroupAPI(qservice.NewWarehouseGroupService())

	fmt.Println("warehouseAPI:", warehouseAPI)
	warehouseRouter := Router.Group("admin/warehouse")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Warehouse true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse [post]
		warehouseRouter.POST("/warehouse", warehouseAPI.CreateWarehouse)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Warehouse
		// @Failure 400 {object} models.Response
		// @Router /warehouse/{id} [get]
		warehouseRouter.GET("/warehouse/:id", warehouseAPI.GetWarehouseByID)

		// @Summary 更新试题
		// @Description 根据ID更新试题信息
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Param body body models.Warehouse true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse/{id} [put]
		warehouseRouter.PUT("/warehouse/:id", warehouseAPI.UpdateWarehouse)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse/{id} [delete]
		warehouseRouter.DELETE("/warehouse/:id", warehouseAPI.DeleteWarehouse)

		// @Summary 批量导入试题
		// @Description 批量导入试题
		// @Tags 试题管理
		// @Accept multipart/form-data
		// @Produce json
		// @Param file formData file true "文件"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse/batch-import [post]
		warehouseRouter.POST("/warehouse/batch-import", warehouseAPI.BatchImportWarehouses)

		// @Summary 导出试题
		// @Description 导出试题
		// @Tags 试题管理
		// @Accept json
		// @Produce octet-stream
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /warehouse/export [get]
		warehouseRouter.GET("/warehouse/export", warehouseAPI.ExportWarehouses)
	}
	{
		// 题组路由
		// @Summary 创建题组
		// @Description 创建一个新的题组
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param body body models.WarehouseGroup true "题组信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse-group [post]
		warehouseRouter.POST("/warehouse-group", warehouseGroupAPI.CreateWarehouseGroup)

		// @Summary 获取题组详情
		// @Description 根据ID获取题组详情
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Success 200 {object} models.WarehouseGroup
		// @Failure 400 {object} models.Response
		// @Router /warehouse-group/{id} [get]
		warehouseRouter.GET("/warehouse-group/:id", warehouseGroupAPI.GetWarehouseGroupByID)

		// @Summary 更新题组
		// @Description 根据ID更新题组信息
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Param body body models.WarehouseGroup true "题组信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse-group/{id} [put]
		warehouseRouter.PUT("/warehouse-group/:id", warehouseGroupAPI.UpdateWarehouseGroup)

		// @Summary 删除题组
		// @Description 根据ID删除题组
		// @Tags 题组管理
		// @Accept json
		// @Produce json
		// @Param id path string true "题组ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /warehouse-group/{id} [delete]
		warehouseRouter.DELETE("/warehouse-group/:id", warehouseGroupAPI.DeleteWarehouseGroup)

		// @Summary 导出题组为PDF
		// @Description 根据ID导出题组为PDF
		// @Tags 题组管理
		// @Accept json
		// @Produce octet-stream
		// @Param id path string true "题组ID"
		// @Success 200 {file} file
		// @Failure 400 {object} models.Response
		// @Router /warehouse-group/export-pdf/{id} [get]
		warehouseRouter.GET("/warehouse-group/export-pdf/:id", warehouseGroupAPI.ExportWarehouseGroupPDF)
	}
}
