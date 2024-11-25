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

type BookAPI struct {
	Service *qservice.BookService
}

func NewBookAPI(service *qservice.BookService) *BookAPI {
	return &BookAPI{Service: service}
}

func (h *BookAPI) CreateBook(c *gin.Context) {
	var q qmodel.RBook
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateBook(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *BookAPI) GetBookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	level := strings.TrimSpace(c.Query("level"))
	majorID := cast.ToInt(c.Query("major_id"))

	books, total, err := h.Service.GetBookList(page, pageSize, keyword, level, majorID)
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

func (h *BookAPI) GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetBookByID(id)
	if err != nil {
		global.GVA_LOG.Error("试题未找到!", zap.Error(err))
		response.FailWithMessage("试题未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *BookAPI) DeleteBook(c *gin.Context) {
	idsStr := c.Param("id")
	ids := parseIDs(idsStr)

	// 检查试题是否存在
	for _, id := range ids {
		existingTopic, err := h.Service.GetBookByID(id)
		if err != nil || existingTopic == nil {
			global.GVA_LOG.Error("试题未找到!", zap.Error(err))
			response.FailWithMessage("部分试题未找到", c)
			return
		}
	}

	if err := h.Service.DeleteBook(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("试题删除成功", c)
}

func parseIDs(idsStr string) []int {
	ids := make([]int, 0)
	for _, idStr := range strings.Split(idsStr, ",") {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			global.GVA_LOG.Error("无效的ID", zap.Error(err))
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
