package initialize

import (
	_ "admin_base_server/source/example"
	_ "admin_base_server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
