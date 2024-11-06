package warehouse

import (
	"admin_base_server/model/common/response"
	qmodel "admin_base_server/model/warehouse"
	qservice "admin_base_server/service/warehouse"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WarehouseGroupAPI struct {
	Service *qservice.WarehouseGroupService
}

func NewWarehouseGroupAPI(service *qservice.WarehouseGroupService) *WarehouseGroupAPI {
	return &WarehouseGroupAPI{Service: service}
}

func (h *WarehouseGroupAPI) CreateWarehouseGroup(c *gin.Context) {
	var g qmodel.WarehouseGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateWarehouseGroup(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *WarehouseGroupAPI) GetWarehouseGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	g, err := h.Service.GetWarehouseGroupByID(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *WarehouseGroupAPI) UpdateWarehouseGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var g qmodel.WarehouseGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.UpdateWarehouseGroup(id, &g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *WarehouseGroupAPI) DeleteWarehouseGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteWarehouseGroup(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("题组删除成功", c)
}

func (h *WarehouseGroupAPI) ExportWarehouseGroupPDF(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pdfData, err := h.Service.ExportWarehouseGroupPDF(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=题组.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfData)))
	c.Writer.Write(pdfData)
}
