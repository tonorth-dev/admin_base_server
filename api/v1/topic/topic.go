package topic

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	qmodel "admin_base_server/model/topic"
	qservice "admin_base_server/service/topic"
	"bytes"
	"encoding/csv"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TopicAPI struct {
	Service *qservice.TopicService
}

func NewTopicAPI(service *qservice.TopicService) *TopicAPI {
	return &TopicAPI{Service: service}
}

func (h *TopicAPI) CreateTopic(c *gin.Context) {
	var q qmodel.Topic
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateTopic(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *TopicAPI) GetTopicList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := strings.TrimSpace(c.Query("keyword"))
	cate := strings.TrimSpace(c.Query("cate"))
	level := strings.TrimSpace(c.Query("level"))
	majorID := cast.ToInt(c.Query("major_id"))

	topics, total, err := h.Service.GetTopicList(page, pageSize, search, cate, level, majorID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  topics,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *TopicAPI) GetTopicByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetTopicByID(id)
	if err != nil {
		global.GVA_LOG.Error("试题未找到!", zap.Error(err))
		response.FailWithMessage("试题未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *TopicAPI) UpdateTopic(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q qmodel.Topic
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查试题是否存在
	existingTopic, err := h.Service.GetTopicByID(id)
	if err != nil || existingTopic == nil {
		global.GVA_LOG.Error("试题未找到!", zap.Error(err))
		response.FailWithMessage("试题未找到", c)
		return
	}

	if err := h.Service.UpdateTopic(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *TopicAPI) DeleteTopic(c *gin.Context) {
	idsStr := c.Param("id")
	ids := parseIDs(idsStr)

	// 检查试题是否存在
	for _, id := range ids {
		existingTopic, err := h.Service.GetTopicByID(id)
		if err != nil || existingTopic == nil {
			global.GVA_LOG.Error("试题未找到!", zap.Error(err))
			response.FailWithMessage("部分试题未找到", c)
			return
		}
	}

	if err := h.Service.DeleteTopic(ids); err != nil {
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

func (h *TopicAPI) BatchImportTopics(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查文件大小
	if file.Size > 20*1024*1024 { // 20MB
		response.FailWithMessage("文件大小超过20MB", c)
		return
	}

	// 读取上传的文件
	fileBytes, err := file.Open()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer fileBytes.Close()

	// 解析 CSV 文件
	reader := csv.NewReader(fileBytes)
	reader.FieldsPerRecord = -1 // 允许不同行有不同的字段数 todo，要关闭
	reader.LazyQuotes = true    // 更宽松地处理引号

	records, err := reader.ReadAll()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查记录数量
	if len(records) > 5000 {
		response.FailWithMessage("内容条数超过5000行", c)
		return
	}

	// 将 CSV 记录转换为 Topic 结构体
	var topics []qmodel.Topic
	for _, record := range records {
		// 防止XSS攻击，对输入数据进行HTML转义
		cate := strings.TrimSpace(html.EscapeString(record[0]))
		level := strings.TrimSpace(html.EscapeString(record[1]))
		title := strings.TrimSpace(html.EscapeString(record[2]))
		answer := strings.TrimSpace(html.EscapeString(record[3]))
		majorID := cast.ToInt(record[4])
		tag := strings.TrimSpace(html.EscapeString(record[5]))
		author := strings.TrimSpace(html.EscapeString(record[6]))

		topic := qmodel.Topic{
			Title:   title,
			Cate:    cate,
			Level:   level,
			Answer:  answer,
			Author:  author,
			MajorID: majorID,
			Tag:     tag,
		}
		topics = append(topics, topic)
	}

	// 分批插入数据
	batchSize := 200
	for i := 0; i < len(topics); i += batchSize {
		end := i + batchSize
		if end > len(topics) {
			end = len(topics)
		}
		batch := topics[i:end]

		if err := h.Service.BatchImportTopics(batch); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	response.OkWithMessage("试题批量导入成功", c)
}

func (h *TopicAPI) ExportTopics(c *gin.Context) {
	// 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := strings.TrimSpace(c.Query("search"))
	level := strings.TrimSpace(c.Query("level"))
	cate := strings.TrimSpace(c.Query("cate"))
	majorID := cast.ToInt(c.Query("major_id"))

	// 调用服务层方法获取符合条件的题目列表
	topics, _, err := h.Service.GetTopicList(page, pageSize, search, cate, level, majorID, []int{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 创建 bytes.Buffer 作为 CSV 写入器的目标
	var csvData bytes.Buffer
	csvWriter := csv.NewWriter(&csvData)

	// 写入 CSV 头
	header := []string{"Title", "Cate", "Answer", "Author", "MajorID", "MajorName", "Tag"}
	if err := csvWriter.Write(header); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 写入试题数据
	for _, topic := range topics {
		record := []string{
			topic.Title,
			cast.ToString(topic.Cate),
			topic.Answer,
			topic.Author,
			cast.ToString(topic.MajorID),
			topic.MajorName,
			topic.Tag,
		}
		if err := csvWriter.Write(record); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	// 完成写入
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=topics.csv")
	c.Header("File-Name", "topics.csv")

	// 发送 CSV 数据
	c.String(http.StatusOK, csvData.String())
}
