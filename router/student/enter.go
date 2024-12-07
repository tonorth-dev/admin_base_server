package student

import (
	"admin_base_server/api/v1/student"
)

type RouterGroup struct {
	StudentRouter
}

var (
	studentAPI *student.StudentAPI
)
