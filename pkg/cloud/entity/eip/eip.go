package eip

import (
	"goutils/pkg/cloud/entity"
)

// Eip 弹性公网IP
type Eip struct {
	entity.Resource `json:"inline" bson:"inline"`
}

// eipData 弹性公网IP数据
type EipData struct {
}

type GetEipRequest struct {
	*entity.BaseCloudRequest
}
