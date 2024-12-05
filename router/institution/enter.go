package job

import (
	"admin_base_server/api/v1/job"
)

type RouterGroup struct {
	JobRouter
}

var (
	jobAPI *job.JobAPI
)
