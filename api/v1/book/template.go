package book

import (
	"admin_base_server/global"
	qmodel "admin_base_server/model/book"
	"admin_base_server/model/common/response"
	qservice "admin_base_server/service/book"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TemplateAPI struct {
	Service *qservice.TemplateService
}

func NewTemplateAPI(service *qservice.TemplateService) *TemplateAPI {
	return &TemplateAPI{Service: service}
}

func (h *TemplateAPI) CreateTemplate(c *gin.Context) {
	var q qmodel.RTemplate
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateTemplate(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *TemplateAPI) GetTemplateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	level := strings.TrimSpace(c.Query("level"))
	majorID := cast.ToInt(c.Query("major_id"))

	books, total, err := h.Service.GetTemplateList(page, pageSize, keyword, level, majorID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  books,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *TemplateAPI) GetTemplateByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetTemplateByID(id)
	if err != nil {
		global.GVA_LOG.Error("模板未找到!", zap.Error(err))
		response.FailWithMessage("模板未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *TemplateAPI) DeleteTemplate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 检查模板是否存在
	existingTemplate, err := h.Service.GetTemplateByID(id)
	if err != nil || existingTemplate == nil {
		global.GVA_LOG.Error("模板未找到!", zap.Error(err))
		response.FailWithMessage("模板未找到", c)
		return
	}

	if err := h.Service.DeleteTemplate(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("模板删除成功", c)
}
