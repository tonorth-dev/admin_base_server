package major

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	qmodel "admin_base_server/model/major"
	qservice "admin_base_server/service/major"
	"encoding/csv"
	"go.uber.org/zap"
	"html"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MajorAPI struct {
	Service *qservice.MajorService
}

func NewMajorAPI(service *qservice.MajorService) *MajorAPI {
	return &MajorAPI{Service: service}
}

func (h *MajorAPI) CreateMajor(c *gin.Context) {
	var q qmodel.Major
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateMajor(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *MajorAPI) GetMajorList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	id, _ := strconv.Atoi(c.Query("major_id"))
	search := strings.TrimSpace(c.Query("keyword"))

	majors, total, err := h.Service.GetMajorList(page, pageSize, id, search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  majors,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *MajorAPI) GetMajorByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetMajorByID(id)
	if err != nil {
		global.GVA_LOG.Error("专业未找到!", zap.Error(err))
		response.FailWithMessage("专业未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *MajorAPI) UpdateMajor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q qmodel.Major
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查专业是否存在
	existingMajor, err := h.Service.GetMajorByID(id)
	if err != nil || existingMajor == nil {
		global.GVA_LOG.Error("专业未找到!", zap.Error(err))
		response.FailWithMessage("专业未找到", c)
		return
	}

	if err := h.Service.UpdateMajor(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *MajorAPI) DeleteMajor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 检查专业是否存在
	existingMajor, err := h.Service.GetMajorByID(id)
	if err != nil || existingMajor == nil {
		global.GVA_LOG.Error("专业未找到!", zap.Error(err))
		response.FailWithMessage("专业未找到", c)
		return
	}

	if err := h.Service.DeleteMajor(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("专业删除成功", c)
}

func (h *MajorAPI) BatchImportMajors(c *gin.Context) {
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
	if len(records) > 2000 {
		response.FailWithMessage("内容条数超过2000行", c)
		return
	}

	// 将 CSV 记录转换为 Major 结构体
	var majors []qmodel.Major
	for _, record := range records {
		// 防止XSS攻击，对输入数据进行HTML转义
		firstLevelCategory := html.EscapeString(record[0])
		secondLevelCategory := html.EscapeString(record[1])
		majorName := html.EscapeString(record[2])
		year := html.EscapeString(record[3])

		major := qmodel.Major{
			FirstLevelCategory:  firstLevelCategory,
			SecondLevelCategory: secondLevelCategory,
			MajorName:           majorName,
			Year:                year,
		}
		majors = append(majors, major)
	}

	// 分批插入数据
	batchSize := 200
	for i := 0; i < len(majors); i += batchSize {
		end := i + batchSize
		if end > len(majors) {
			end = len(majors)
		}
		batch := majors[i:end]

		if err := h.Service.BatchImportMajors(batch); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}

	response.OkWithMessage("专业批量导入成功", c)
}

//func (h *MajorAPI) ExportMajors(c *gin.Context) {
//	// 获取查询参数
//	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
//	search := strings.TrimSpace(c.Query("search"))
//	cateStr := strings.TrimSpace(c.Query("cate"))
//	majorIDStr := strings.TrimSpace(c.Query("major_id"))
//
//	// 将字符串参数转换为整数
//	cate := cast.ToInt(cateStr)
//	majorID := cast.ToInt(majorIDStr)
//
//	// 调用服务层方法获取符合条件的题目列表
//	majors, _, err := h.Service.GetMajorList(page, pageSize, search)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//
//	// 创建 bytes.Buffer 作为 CSV 写入器的目标
//	var csvData bytes.Buffer
//	csvWriter := csv.NewWriter(&csvData)
//
//	// 写入 CSV 头
//	header := []string{"Title", "Cate", "Answer", "Author", "MajorID", "MajorName", "Tag"}
//	if err := csvWriter.Write(header); err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//
//	// 写入专业数据
//	for _, major := range majors {
//		record := []string{
//			major.Title,
//			cast.ToString(major.Cate),
//			major.Answer,
//			major.Author,
//			cast.ToString(major.MajorID),
//			major.MajorName,
//			major.Tag,
//		}
//		if err := csvWriter.Write(record); err != nil {
//			response.FailWithMessage(err.Error(), c)
//			return
//		}
//	}
//
//	// 完成写入
//	csvWriter.Flush()
//	if err := csvWriter.Error(); err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//
//	// 设置响应头
//	c.Header("Content-Type", "text/csv")
//	c.Header("Content-Disposition", "attachment; filename=majors.csv")
//	c.Header("File-Name", "majors.csv")
//
//	// 发送 CSV 数据
//	c.String(http.StatusOK, csvData.String())
//}
