package question

import (
	qmodel "admin_base_server/model/question"
	qservice "admin_base_server/service/question"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionGroupAPI struct {
	Service *qservice.QuestionGroupService
}

func NewQuestionGroupAPI(service *qservice.QuestionGroupService) *QuestionGroupAPI {
	return &QuestionGroupAPI{Service: service}
}

func (h *QuestionGroupAPI) CreateQuestionGroup(c *gin.Context) {
	var g qmodel.QuestionGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateQuestionGroup(&g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *QuestionGroupAPI) GetQuestionGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	g, err := h.Service.GetQuestionGroupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "题组未找到"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *QuestionGroupAPI) UpdateQuestionGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var g qmodel.QuestionGroup
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateQuestionGroup(id, &g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *QuestionGroupAPI) DeleteQuestionGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteQuestionGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "题组删除成功"})
}

func (h *QuestionGroupAPI) ExportQuestionGroupPDF(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	pdfData, err := h.Service.ExportQuestionGroupPDF(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=题组.pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfData)))
	c.Writer.Write(pdfData)
}
