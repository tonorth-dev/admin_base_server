package router

import (
	"admin_base_server/router/example"
	"admin_base_server/router/question"
	"admin_base_server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Question question.QuestionRouter
}
