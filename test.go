package main

import (
	"fmt"
	"github.com/signintech/gopdf"
	"log"
	"os"
	"path/filepath"
)

// 定义数据结构
type Question struct {
	ID     int
	Title  string
	Answer string
}

type Section struct {
	Title           string
	QuestionsDetail []Category
}

type Category struct {
	CateName string
	List     []Question
}

type Data struct {
	Name            string
	MajorName       string
	LevelName       string
	UnitNumber      int
	ComponentDesc   []string
	QuestionsNumber int
	QuestionsDesc   []Section
}

// 函数来创建PDF
func exportPDF(isTeacherVersion bool, data Data) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err := pdf.AddTTFFont("OPPOSans-Regular", "static/fonts/OPPOSans-Regular.ttf")
	if err != nil {
		return fmt.Errorf("failed to add font: %v", err)
	}
	pdf.SetFont("OPPOSans-Regular", "", 14)

	// 添加页面
	pdf.AddPage()

	// 标题
	pdf.SetFont("OPPOSans-Regular", "B", 24)
	pdf.SetTextColor(0, 0, 255) // 蓝色
	pdf.Cell(nil, fmt.Sprintf("%s", data.Name))
	pdf.Br(20)

	// 表头信息
	pdf.SetFont("OPPOSans-Regular", "", 12)
	pdf.SetTextColor(0, 0, 255) // 蓝色
	pdf.Cell(nil, fmt.Sprintf("专业：%s", data.MajorName))
	pdf.Br(5)
	pdf.Cell(nil, fmt.Sprintf("难度：%s", data.LevelName))
	pdf.Br(5)
	pdf.Cell(nil, fmt.Sprintf("试题套数：%d", data.UnitNumber))
	pdf.Br(5)

	// 试题组成和总数
	pdf.Cell(nil, fmt.Sprintf("试题组成：%s", joinStrings(data.ComponentDesc, "，")))
	pdf.Br(5)
	pdf.Cell(nil, fmt.Sprintf("试题总数：%d", data.QuestionsNumber))
	pdf.Br(10)

	// 添加试题部分
	for _, section := range data.QuestionsDesc {
		pdf.SetFont("OPPOSans-Regular", "B", 18)
		pdf.SetTextColor(0, 0, 255) // 蓝色
		pdf.Cell(nil, fmt.Sprintf("章节：%s", section.Title))
		pdf.Br(10)

		for _, detail := range section.QuestionsDetail {
			pdf.SetFont("OPPOSans-Regular", "B", 14)
			pdf.SetTextColor(0, 0, 255) // 蓝色
			pdf.Cell(nil, fmt.Sprintf("试题分类：%s", detail.CateName))
			pdf.Br(5)

			for _, question := range detail.List {
				pdf.SetFont("OPPOSans-Regular", "", 12)
				pdf.Cell(nil, fmt.Sprintf("序号: %d", question.ID))
				pdf.Br(5)
				pdf.Cell(nil, fmt.Sprintf("试题标题: %s", question.Title))
				pdf.Br(5)
				pdf.Cell(nil, fmt.Sprintf("试题答案: %s", question.Answer))
				pdf.Br(5)
				pdf.Line(10, pdf.GetY(), 580, pdf.GetY()) // 分割线
				pdf.Br(10)
			}
		}
	}

	// 保存PDF文件
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	fileName := "学生版题本.pdf"
	if isTeacherVersion {
		fileName = "教师版题本.pdf"
	}

	filePath := filepath.Join(dir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create PDF file: %v", err)
	}
	defer file.Close()

	err = pdf.WritePdf(file.Name())
	if err != nil {
		return fmt.Errorf("failed to write PDF: %v", err)
	}

	fmt.Printf("PDF已保存到：%s\n", filePath)
	return nil
}

// 辅助函数，用于连接字符串
func joinStrings(s []string, separator string) string {
	if len(s) == 0 {
		return "无"
	}
	return s[0] + separator + joinStrings(s[1:], separator)
}

func main1() {
	// 这里可以添加你的测试数据
	data := Data{
		Name:            "题本详情",
		MajorName:       "计算机科学",
		LevelName:       "中级",
		UnitNumber:      10,
		ComponentDesc:   []string{"选择题", "填空题", "判断题"},
		QuestionsNumber: 50,
		QuestionsDesc: []Section{
			{
				Title: "第一章",
				QuestionsDetail: []Category{
					{
						CateName: "选择题",
						List: []Question{
							{ID: 1, Title: "这是一个选择题", Answer: "A"},
						},
					},
				},
			},
		},
	}

	err := exportPDF(true, data)
	if err != nil {
		log.Fatalf("PDF导出时发生错误: %v", err)
	}
}
