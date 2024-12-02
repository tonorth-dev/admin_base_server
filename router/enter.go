package router

import (
	"admin_base_server/router/book"
	"admin_base_server/router/config"
	"admin_base_server/router/example"
	"admin_base_server/router/job"
	"admin_base_server/router/major"
	"admin_base_server/router/system"
	"admin_base_server/router/topic"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Topic   topic.RouterGroup
	Major   major.RouterGroup
	Config  config.RouterGroup
	Book    book.RouterGroup
	Job     job.RouterGroup
}
