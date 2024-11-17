package book

import (
	"admin_base_server/api/v1/book"
)

type RouterGroup struct {
	BookRouter
}

var (
	bookAPI *book.BookAPI
)
