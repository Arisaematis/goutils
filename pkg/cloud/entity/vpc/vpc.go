package vpc

import (
	"goutils/pkg/cloud/entity"
)

// Vpc 虚拟私有云
type Vpc struct {
	entity.Resource `json:"inline" bson:"inline"`
}

// VpcData 虚拟私有云数据
type VpcData struct {
	CidrBlock string `json:"cidrBlock" bson:"cidrBlock"`
}

// GetVpcRequest 获取虚拟私有云
type GetVpcRequest struct {
	*entity.BaseCloudRequest
}
