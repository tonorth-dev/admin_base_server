package warehouse

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	qmodel "admin_base_server/model/warehouse"
	qservice "admin_base_server/service/warehouse"
	"encoding/json"
	"go.uber.org/zap"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WarehouseAPI struct {
	Service *qservice.WarehouseService
}

func NewWarehouseAPI(service *qservice.WarehouseService) *WarehouseAPI {
	return &WarehouseAPI{Service: service}
}

func (h *WarehouseAPI) CreateWarehouse(c *gin.Context) {
	var q qmodel.Warehouse
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateWarehouse(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *WarehouseAPI) GetWarehouseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetWarehouseByID(id)
	if err != nil {
		global.GVA_LOG.Error("试题未找到!", zap.Error(err))
		response.FailWithMessage("试题未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *WarehouseAPI) UpdateWarehouse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q qmodel.Warehouse
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.UpdateWarehouse(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *WarehouseAPI) DeleteWarehouse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteWarehouse(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("试题删除成功", c)
}

func (h *WarehouseAPI) BatchImportWarehouses(c *gin.Context) {
	body := c.Request.Body
	var warehouses []qmodel.Warehouse
	if err := json.NewDecoder(body).Decode(&warehouses); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.BatchImportWarehouses(warehouses); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("试题批量导入成功", c)
}

func (h *WarehouseAPI) ExportWarehouses(c *gin.Context) {
	warehouses, err := h.Service.ExportWarehouses()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(warehouses, c)
}
