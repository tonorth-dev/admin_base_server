package institution

import (
	"admin_base_server/api/v1/institution"
)

type RouterGroup struct {
	InstitutionRouter
}

var (
	institutionAPI *institution.InstitutionAPI
)
