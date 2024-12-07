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
	configs, err := h.Service.GetAllDBConfigList()
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
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		response.FailWithMessage("参数name不能为空", c)
		return
	}
	if name == "area" {
		level := strings.TrimSpace(c.Query("level"))
		parentId := strings.TrimSpace(c.Query("parent_id"))
		q, err := h.Service.GetAreaConfig(level, parentId)
		if err != nil {
			global.GVA_LOG.Error("配置未找到!", zap.Error(err))
			response.FailWithMessage("配置未找到", c)
			return
		}
		response.OkWithData(q, c)
		return
	}

	q, err := h.Service.GetActiveConfigFromDB(name)
	if err != nil {
		global.GVA_LOG.Error("配置未找到!", zap.Error(err))
		response.FailWithMessage("配置未找到", c)
		return
	}
	response.OkWithData(q, c)
}

func (h *ConfigAPI) UpdateConfig(c *gin.Context) {
	name := strings.TrimSpace(c.Param("name"))
	var q cmodel.RConfig
	if err := c.ShouldBindJSON(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查配置是否存在
	existingConfig, err := h.Service.GetActiveConfigFromDB(name)
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
