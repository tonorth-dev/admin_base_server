package v1

import (
	"admin_base_server/api/v1/example"
	"admin_base_server/api/v1/system"
	"admin_base_server/api/v1/warehouse"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	QuestionApiGroup warehouse.ApiGroup
}
