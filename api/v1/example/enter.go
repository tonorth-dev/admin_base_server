package example

import "admin_base_server/service"

type ApiGroup struct {
	CustomerApi
	FileUploadAndDownloadApi
}

var (
	customerService              = service.ServiceGroupApp.ExampleServiceGroup.CustomerService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
