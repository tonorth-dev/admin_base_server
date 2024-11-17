package book

import (
	"admin_base_server/api/v1/book"
	qservice "admin_base_server/service/book"
	"github.com/gin-gonic/gin"
)

// BookRouter 路由管理器
type BookRouter struct{}

// @Summary 初始化试题和题组路由
// @Description 初始化试题和题组的相关路由
// @Tags admin
// @Router /admin/book [post,get,put,delete]
// @Router /admin/book-group [post,get,put,delete]
func (e *BookRouter) InitBookRouter(Router *gin.RouterGroup) {

	bookAPI = book.NewBookAPI(qservice.NewBookService())

	bookRouter := Router.Group("admin/book")
	{
		// 试题路由
		// @Summary 创建试题
		// @Description 创建一个新的试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param body body models.Book true "试题信息"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /book [post]
		bookRouter.POST("/book", bookAPI.CreateBook)

		// @Summary 获取试题详情
		// @Description 根据ID获取试题详情
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Book
		// @Failure 400 {object} models.Response
		// @Router /book/{id} [get]
		bookRouter.GET("/book/:id", bookAPI.GetBookByID)

		// @Summary 获取所有试题列表
		// @Description 获取所有试题列表
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Success 200 {array} models.Book
		// @Failure 400 {object} models.Response
		// @Router /book/list [get]
		bookRouter.GET("/book/list", bookAPI.GetBookList)

		// @Summary 删除试题
		// @Description 根据ID删除试题
		// @Tags 试题管理
		// @Accept json
		// @Produce json
		// @Param id path string true "试题ID"
		// @Success 200 {object} models.Response
		// @Failure 400 {object} models.Response
		// @Router /book/{id} [delete]
		bookRouter.DELETE("/book/:id", bookAPI.DeleteBook)
	}
}
