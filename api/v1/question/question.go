package question

import (
	qmodel "admin_base_server/model/question"
	qservice "admin_base_server/service/question"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionAPI struct {
	Service *qservice.QuestionService
}

func NewQuestionAPI(service *qservice.QuestionService) *QuestionAPI {
	return &QuestionAPI{Service: service}
}

func (h *QuestionAPI) CreateQuestion(c *gin.Context) {
	var q qmodel.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateQuestion(&q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (h *QuestionAPI) GetQuestionByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "试题未找到"})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (h *QuestionAPI) UpdateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q qmodel.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.UpdateQuestion(id, &q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (h *QuestionAPI) DeleteQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Service.DeleteQuestion(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "试题删除成功"})
}

func (h *QuestionAPI) BatchImportQuestions(c *gin.Context) {
	body := c.Request.Body
	var questions []qmodel.Question
	if err := json.NewDecoder(body).Decode(&questions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.BatchImportQuestions(questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "试题批量导入成功"})
}

func (h *QuestionAPI) ExportQuestions(c *gin.Context) {
	questions, err := h.Service.ExportQuestions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}
