package cloud_request

import (
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	accessKey "goutils/pkg/cloud/entity/accesskey"
)

type BaseRequestProducer struct {
	accessKey.AccessKey
	serviceType.Type
}
