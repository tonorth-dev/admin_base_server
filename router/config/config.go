package config

import (
	"admin_base_server/api/v1/config"
	qservice "admin_base_server/service/config"
	"github.com/gin-gonic/gin"
)

// ConfigRouter 路由管理器
type ConfigRouter struct{}

// @Summary 初始化配置和题组路由
// @Description 初始化配置和题组的相关路由
// @Tags admin
// @Router /admin/config [post,get,put,delete]
func (e *ConfigRouter) InitConfigRouter(Router *gin.RouterGroup) {
	configAPI = config.NewConfigAPI(qservice.NewConfigService())

	configGroup := Router.Group("admin/config")
	{
		configGroup.POST("/config/", configAPI.CreateConfig)
		configGroup.GET("/config/list", configAPI.GetAllConfigList)
		configGroup.GET("config/:name", configAPI.GetConfigByName)
		configGroup.PUT("config/:name", configAPI.UpdateConfig)
	}
}
