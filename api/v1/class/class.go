package class

import (
	"admin_base_server/global"
	jmodel "admin_base_server/model/class"
	"admin_base_server/model/common/response"
	jservice "admin_base_server/service/class"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type ClassAPI struct {
	Service *jservice.ClassService
}

func NewClassAPI(service *jservice.ClassService) *ClassAPI {
	return &ClassAPI{Service: service}
}

func (h *ClassAPI) CreateClass(c *gin.Context) {
	var q jmodel.Class
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateClass(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *ClassAPI) GetClassList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	institutionId := cast.ToInt(c.Query("institution_id"))

	var classs []jmodel.RClass
	var total int64
	var err error

	classs, total, err = h.Service.GetClassList(page, pageSize, keyword, institutionId)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  classs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *ClassAPI) GetClassByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetClassByID(id)
	if err != nil {
		global.GVA_LOG.Error("班级未找到!", zap.Error(err))
		response.FailWithMessage("班级未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *ClassAPI) UpdateClass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q jmodel.Class
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查班级是否存在
	existingClass, err := h.Service.GetClassByID(id)
	if err != nil || existingClass == nil {
		global.GVA_LOG.Error("班级未找到!", zap.Error(err))
		response.FailWithMessage("班级未找到", c)
		return
	}

	if err := h.Service.UpdateClass(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("更新成功", c)
}

func (h *ClassAPI) DeleteClass(c *gin.Context) {
	idsStr := c.Param("id")
	ids := parseIDs(idsStr)

	// 检查班级是否存在
	for _, id := range ids {
		existingClass, err := h.Service.GetClassByID(id)
		if err != nil || existingClass == nil {
			global.GVA_LOG.Error("班级未找到!", zap.Error(err))
			response.FailWithMessage("部分班级未找到", c)
			return
		}
	}

	if err := h.Service.DeleteClass(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("班级删除成功", c)
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

// BatchImportClasss 用于处理批量导入题目，要根据文件格式和内容调整
func (h *ClassAPI) BatchImportClasss(ctx *gin.Context) {
	// 获取上传的文件
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	// 打开文件以获取 io.Reader
	file, err := fileHeader.Open()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	defer file.Close()

	// 读取Excel文件数据
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	defer xlsx.Close()

	// 获取指定工作表中的数据
	var records [][]string
	sheetName := xlsx.GetSheetName(0) // 获取第一个工作表名称
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	records = append(records, rows...)

	// 检查记录数量
	if len(records) > 5000 {
		response.FailWithMessage("内容条数超过5000行", ctx)
		return
	}

	// 将解析的数据转为Class结构体

	response.OkWithMessage("班级批量导入成功", ctx)
}

func findID(data []map[string]string, value string) string {
	for _, item := range data {
		if item["name"] == value {
			return item["id"]
		}
	}
	return ""
}

func (h *ClassAPI) ExportClasss(c *gin.Context) {
	//// 获取查询参数
	//page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	//pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	//search := strings.TrimSpace(c.Query("search"))
	//level := strings.TrimSpace(c.Query("level"))
	//cate := strings.TrimSpace(c.Query("cate"))
	//majorID := cast.ToInt(c.Query("major_id"))
	//status := cast.ToInt(c.Query("status"))
	//
	//// 调用服务层方法获取符合条件的题目列表
	//classs, _, err := h.Service.GetClassList(page, pageSize, search, cate, level, majorID, status)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//// 创建 bytes.Buffer 作为 CSV 写入器的目标
	//var csvData bytes.Buffer
	//csvWriter := csv.NewWriter(&csvData)
	//
	//// 写入 CSV 头
	//header := []string{"Title", "Cate", "Answer", "Author", "MajorID", "MajorName", "Tag"}
	//if err := csvWriter.Write(header); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//// 写入班级数据
	//for _, class := range classs {
	//	record := []string{
	//		class.Title,
	//		cast.ToString(class.Cate),
	//		class.Answer,
	//		class.Author,
	//		cast.ToString(class.MajorID),
	//		class.MajorName,
	//		class.Tag,
	//	}
	//	if err := csvWriter.Write(record); err != nil {
	//		response.FailWithMessage(err.Error(), c)
	//		return
	//	}
	//}
	//
	//// 完成写入
	//csvWriter.Flush()
	//if err := csvWriter.Error(); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//// 设置响应头
	//c.Header("Content-Type", "text/csv")
	//c.Header("Content-Disposition", "attachment; filename=classs.csv")
	//c.Header("File-Name", "classs.csv")

	// 发送 CSV 数据
	//c.String(http.StatusOK, csvData.String())
}
