package config

import (
	"admin_base_server/api/v1/config"
)

type RouterGroup struct {
	ConfigRouter
}

var (
	configAPI *config.ConfigAPI
)
