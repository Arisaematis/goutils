package rds

import (
	"goutils/pkg/cloud/entity"
)

// Rds 云数据库
type Rds struct {
	entity.Resource `json:"inline" bson:"inline"`
}

type RdsData struct {
	PrivateIps      []string `json:"privateIps" bson:"privateIps"`
	PrivateDnsNames []string `json:"privateDnsNames" bson:"privateDnsNames"`
	PublicIps       []string `json:"publicIps" bson:"publicIps"`
	Port            int64    `json:"port" bson:"port"`
}

// GetRdsRequest 获取云数据库
type GetRdsRequest struct {
	*entity.BaseCloudRequest
}
