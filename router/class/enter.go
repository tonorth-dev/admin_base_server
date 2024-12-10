package class

import (
	"admin_base_server/api/v1/class"
)

type RouterGroup struct {
	ClassRouter
}

var (
	classAPI *class.ClassAPI
)
