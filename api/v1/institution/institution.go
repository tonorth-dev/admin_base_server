package institution

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	jmodel "admin_base_server/model/institution"
	jservice "admin_base_server/service/institution"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type InstitutionAPI struct {
	Service *jservice.InstitutionService
}

func NewInstitutionAPI(service *jservice.InstitutionService) *InstitutionAPI {
	return &InstitutionAPI{Service: service}
}

func (h *InstitutionAPI) CreateInstitution(c *gin.Context) {
	var q jmodel.Institution
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateInstitution(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *InstitutionAPI) GetInstitutionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	province := strings.TrimSpace(c.Query("province"))
	city := strings.TrimSpace(c.Query("city"))

	var institutions []jmodel.RInstitution
	var total int64
	var err error

	institutions, total, err = h.Service.GetInstitutionList(page, pageSize, keyword, province, city)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  institutions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *InstitutionAPI) GetInstitutionByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetInstitutionByID(id)
	if err != nil {
		global.GVA_LOG.Error("机构未找到!", zap.Error(err))
		response.FailWithMessage("机构未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *InstitutionAPI) UpdateInstitution(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q jmodel.Institution
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查机构是否存在
	existingInstitution, err := h.Service.GetInstitutionByID(id)
	if err != nil || existingInstitution == nil {
		global.GVA_LOG.Error("机构未找到!", zap.Error(err))
		response.FailWithMessage("机构未找到", c)
		return
	}

	if err := h.Service.UpdateInstitution(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("更新成功", c)
}

func (h *InstitutionAPI) DeleteInstitution(c *gin.Context) {
	idsStr := c.Param("id")
	ids := parseIDs(idsStr)

	// 检查机构是否存在
	for _, id := range ids {
		existingInstitution, err := h.Service.GetInstitutionByID(id)
		if err != nil || existingInstitution == nil {
			global.GVA_LOG.Error("机构未找到!", zap.Error(err))
			response.FailWithMessage("部分机构未找到", c)
			return
		}
	}

	if err := h.Service.DeleteInstitution(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("机构删除成功", c)
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

// BatchImportInstitutions 用于处理批量导入题目，要根据文件格式和内容调整
func (h *InstitutionAPI) BatchImportInstitutions(ctx *gin.Context) {
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

	response.OkWithMessage("机构批量导入成功", ctx)
}

func findID(data []map[string]string, value string) string {
	for _, item := range data {
		if item["name"] == value {
			return item["id"]
		}
	}
	return ""
}

func (h *InstitutionAPI) ExportInstitutions(c *gin.Context) {
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
	//institutions, _, err := h.Service.GetInstitutionList(page, pageSize, search, cate, level, majorID, status)
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
	//// 写入机构数据
	//for _, institution := range institutions {
	//	record := []string{
	//		institution.Title,
	//		cast.ToString(institution.Cate),
	//		institution.Answer,
	//		institution.Author,
	//		cast.ToString(institution.MajorID),
	//		institution.MajorName,
	//		institution.Tag,
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
	//c.Header("Content-Disposition", "attachment; filename=institutions.csv")
	//c.Header("File-Name", "institutions.csv")

	// 发送 CSV 数据
	//c.String(http.StatusOK, csvData.String())
}
