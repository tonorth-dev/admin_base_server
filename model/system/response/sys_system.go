package response

import "admin_base_server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
