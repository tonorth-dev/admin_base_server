package major

import (
	"admin_base_server/api/v1/major"
)

type RouterGroup struct {
	MajorRouter
}

var (
	majorAPI *major.MajorAPI
)
