package request

import (
	"admin_base_server/model/common/request"
	"admin_base_server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
