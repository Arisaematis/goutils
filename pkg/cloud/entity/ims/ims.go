package ims

import (
	"goutils/pkg/cloud/entity"
)

// Ims 镜像服务
type Ims struct {
	entity.Resource `json:"inline" bson:"inline"`
}

type ImsData struct {
}

// GetImsRequest 获取镜像服务
type GetImsRequest struct {
	*entity.BaseCloudRequest
}
