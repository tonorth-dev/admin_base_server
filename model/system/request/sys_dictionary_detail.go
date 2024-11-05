package request

import (
	"admin_base_server/model/common/request"
	"admin_base_server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
