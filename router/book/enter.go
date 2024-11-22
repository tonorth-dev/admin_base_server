package book

import (
	"admin_base_server/api/v1/book"
)

type RouterGroup struct {
	BookRouter
	TemplateRouter
}

var (
	bookAPI     *book.BookAPI
	templateAPI *book.TemplateAPI
)
