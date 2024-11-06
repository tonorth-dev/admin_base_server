package warehouse

import (
	"admin_base_server/api/v1/warehouse"
)

type RouterGroup struct {
	WarehouseRouter
}

var (
	warehouseAPI      *warehouse.WarehouseAPI
	warehouseGroupAPI *warehouse.WarehouseGroupAPI
)
