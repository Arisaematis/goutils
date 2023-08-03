package obs

import (
	"goutils/pkg/cloud/entity"
	"goutils/pkg/cloud/entity/request"
)

type GetObsArgs struct {
	*entity.BaseCloudRequest
}

type Bucket struct {
	request.BaseCloudResponse
	BucketType string  `json:"bucketType,omitempty"`
	Size       float64 `json:"size"`
	QuotaSize  float64 `json:"quotaSize"`
}
