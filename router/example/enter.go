package example

import (
	api "admin_base_server/api/v1"
)

type RouterGroup struct {
	CustomerRouter
	FileUploadAndDownloadRouter
}

var (
	exaCustomerApi              = api.ApiGroupApp.ExampleApiGroup.CustomerApi
	exaFileUploadAndDownloadApi = api.ApiGroupApp.ExampleApiGroup.FileUploadAndDownloadApi
)
