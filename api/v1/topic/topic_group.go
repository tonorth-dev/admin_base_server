package topic

import (
	"admin_base_server/model/common/response"
	qmodel "admin_base_server/model/topic"
	qservice "admin_base_server/service/topic"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TopicGroupAPI struct {
	Service *qservice.TopicGroupService
}

func NewTopicGroupAPI(service *qservice.TopicGroupService) *TopicGroupAPI {
	return &TopicGroupAPI{Service: service}
}

func (h *TopicGroupAPI) CreateTopicGroup(c *gin.Context) {
	var g qmodel.TopicGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateTopicGroup(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *TopicGroupAPI) GetTopicGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	g, err := h.Service.GetTopicGroupByID(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *TopicGroupAPI) UpdateTopicGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var g qmodel.TopicGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.UpdateTopicGroup(id, &g); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(g, c)
}

func (h *TopicGroupAPI) DeleteTopicGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteTopicGroup(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("题组删除成功", c)
}

func (h *TopicGroupAPI) ExportTopicGroupPDF(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pdfData, err := h.Service.ExportTopicGroupPDF(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=题组.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfData)))
	c.Writer.Write(pdfData)
}
