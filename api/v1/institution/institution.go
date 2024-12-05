package job

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	jmodel "admin_base_server/model/job"
	jservice "admin_base_server/service/job"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"html"
	"strconv"
	"strings"
)

type JobAPI struct {
	Service *jservice.JobService
}

func NewJobAPI(service *jservice.JobService) *JobAPI {
	return &JobAPI{Service: service}
}

func (h *JobAPI) CreateJob(c *gin.Context) {
	var q jmodel.Job
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.CreateJob(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *JobAPI) GetJobList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	search := strings.TrimSpace(c.Query("keyword"))
	majorID := cast.ToInt(c.Query("major_id"))
	all := cast.ToInt(c.Query("all"))

	var jobs []jmodel.RJob
	var total int64
	var err error

	if majorID > 0 && all > 0 {
		jobs, total, err = h.Service.GetJobListBySortMajor(page, pageSize, search, majorID)
	} else {
		jobs, total, err = h.Service.GetJobList(page, pageSize, search, majorID)
	}

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  jobs,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, c)
}

func (h *JobAPI) GetJobByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q, err := h.Service.GetJobByID(id)
	if err != nil {
		global.GVA_LOG.Error("岗位未找到!", zap.Error(err))
		response.FailWithMessage("岗位未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *JobAPI) UpdateJob(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var q jmodel.Job
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查岗位是否存在
	existingJob, err := h.Service.GetJobByID(id)
	if err != nil || existingJob == nil {
		global.GVA_LOG.Error("岗位未找到!", zap.Error(err))
		response.FailWithMessage("岗位未找到", c)
		return
	}

	if err := h.Service.UpdateJob(id, &q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("更新成功", c)
}

func (h *JobAPI) BatchUpdateMajor(c *gin.Context) {
	var req jmodel.RBatchUpdateMajor
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 批量更新 major_id
	if err := h.Service.BatchUpdateMajor(req.JobIDs, req.MajorID); err != nil {
		global.GVA_LOG.Error("批量更新 major_id 失败!", zap.Error(err))
		response.FailWithMessage("批量更新 major_id 失败", c)
		return
	}

	response.OkWithData("批量更新成功", c)
}

func (h *JobAPI) DeleteJob(c *gin.Context) {
	idsStr := c.Param("id")
	ids := parseIDs(idsStr)

	// 检查岗位是否存在
	for _, id := range ids {
		existingJob, err := h.Service.GetJobByID(id)
		if err != nil || existingJob == nil {
			global.GVA_LOG.Error("岗位未找到!", zap.Error(err))
			response.FailWithMessage("部分岗位未找到", c)
			return
		}
	}

	if err := h.Service.DeleteJob(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("岗位删除成功", c)
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

// BatchImportJobs 用于处理批量导入题目，要根据文件格式和内容调整
func (h *JobAPI) BatchImportJobs(ctx *gin.Context) {
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

	// 将解析的数据转为Job结构体
	var jobs []jmodel.Job
	for _, record := range records[5:] { // 跳过标题行
		if len(record) < 20 { // 根据Job结构体中字段数量判断
			fmt.Printf("数据行字段数量不足: %v\n", record)
			continue
		}

		// 防止XSS攻击，对输入数据进行HTML转义
		code := strings.TrimSpace(html.EscapeString(record[0]))
		name := strings.TrimSpace(html.EscapeString(record[4]))
		cate := strings.TrimSpace(html.EscapeString(record[3]))
		desc := strings.TrimSpace(html.EscapeString(record[5]))
		companyCode := strings.TrimSpace(html.EscapeString(record[1]))
		companyName := strings.TrimSpace(html.EscapeString(record[2]))
		enrollmentNumStr := strings.TrimSpace(html.EscapeString(record[6]))
		enrollmentRatio := strings.TrimSpace(html.EscapeString(record[7]))
		city := strings.TrimSpace(html.EscapeString(record[18]))
		phone := strings.TrimSpace(html.EscapeString(record[19]))

		// 解析整数字段
		enrollmentNum, err := strconv.Atoi(enrollmentNumStr)
		if err != nil {
			fmt.Printf("转换EnrollmentNum失败: %v\n", record)
			continue
		}

		// 构造其他条件
		condition := map[string]string{
			"exam":          strings.TrimSpace(html.EscapeString(record[12])),
			"major":         strings.TrimSpace(html.EscapeString(record[11])),
			"other":         strings.TrimSpace(html.EscapeString(record[17])),
			"source":        findID(jmodel.Source, strings.TrimSpace(record[8])),
			"qualification": findID(jmodel.Qualification, strings.TrimSpace(record[9])),
			"degree":        findID(jmodel.Degree, strings.TrimSpace(record[10])),
		}
		conditionJson, _ := json.Marshal(condition)

		job := jmodel.Job{
			Code:            code,
			Name:            name,
			Cate:            cate,
			Desc:            desc,
			CompanyCode:     companyCode,
			CompanyName:     companyName,
			EnrollmentNum:   enrollmentNum,
			EnrollmentRatio: enrollmentRatio,
			Condition:       conditionJson,
			City:            city,
			Phone:           phone,
			Status:          2, // 假设状态为已完成
		}
		jobs = append(jobs, job)
	}

	// 使用事务进行批量插入以确保数据一致性
	// 分批插入数据
	batchSize := 200
	for i := 0; i < len(jobs); i += batchSize {
		end := i + batchSize
		if end > len(jobs) {
			end = len(jobs)
		}
		batch := jobs[i:end]

		if err := h.Service.BatchImportJobs(batch); err != nil {
			response.FailWithMessage(err.Error(), ctx)
			return
		}
	}

	response.OkWithMessage("岗位批量导入成功", ctx)
}

func findID(data []map[string]string, value string) string {
	for _, item := range data {
		if item["name"] == value {
			return item["id"]
		}
	}
	return ""
}

func (h *JobAPI) ExportJobs(c *gin.Context) {
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
	//jobs, _, err := h.Service.GetJobList(page, pageSize, search, cate, level, majorID, status)
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
	//// 写入岗位数据
	//for _, job := range jobs {
	//	record := []string{
	//		job.Title,
	//		cast.ToString(job.Cate),
	//		job.Answer,
	//		job.Author,
	//		cast.ToString(job.MajorID),
	//		job.MajorName,
	//		job.Tag,
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
	//c.Header("Content-Disposition", "attachment; filename=jobs.csv")
	//c.Header("File-Name", "jobs.csv")

	// 发送 CSV 数据
	//c.String(http.StatusOK, csvData.String())
}
