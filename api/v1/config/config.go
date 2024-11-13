package config

import (
	"admin_base_server/global"
	"admin_base_server/model/common/response"
	cmodel "admin_base_server/model/config"
	cservice "admin_base_server/service/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

type ConfigAPI struct {
	Service *cservice.ConfigService
}

func NewConfigAPI(service *cservice.ConfigService) *ConfigAPI {
	return &ConfigAPI{Service: service}
}

func (h *ConfigAPI) CreateConfig(c *gin.Context) {
	var q cmodel.RConfig
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := h.Service.InsertConfig(q.Name, q.Attr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}

func (h *ConfigAPI) GetAllConfigList(c *gin.Context) {
	configs, err := h.Service.GetAllConfigList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"list":  configs,
		"total": len(configs),
	}, c)
}

func (h *ConfigAPI) GetConfigByName(c *gin.Context) {
	name := strings.Trim(c.Param("name"), " ")
	q, err := h.Service.GetActiveConfig(name)
	if err != nil {
		global.GVA_LOG.Error("配置未找到!", zap.Error(err))
		response.FailWithMessage("配置未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *ConfigAPI) UpdateConfig(c *gin.Context) {
	name := strings.Trim(c.Param("name"), " ")
	var q cmodel.RConfig
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查配置是否存在
	existingConfig, err := h.Service.GetActiveConfig(name)
	if err != nil || existingConfig == nil {
		global.GVA_LOG.Error("配置未找到!", zap.Error(err))
		response.FailWithMessage("配置未找到", c)
		return
	}

	if err := h.Service.UpdateConfig(name, q.Attr); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(q, c)
}
