package router

import (
	"admin_base_server/router/example"
	"admin_base_server/router/system"
	"admin_base_server/router/warehouse"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System    system.RouterGroup
	Example   example.RouterGroup
	Warehouse warehouse.RouterGroup
}
