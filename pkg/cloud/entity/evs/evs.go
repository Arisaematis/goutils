package evs

import (
	"goutils/pkg/cloud/entity"
)

type Disk struct {
	Id          string `json:"id" bson:"id"`
	Size        int    `json:"size" bson:"size"`
	Category    string `json:"category" bson:"category"`
	DiskType    string `json:"diskType" bson:"diskType"` // system/data
	DiskName    string `json:"diskName" bson:"diskName"`
	InstanceId  string `json:"instanceId" bson:"instanceId"`
	Description string `json:"description" bson:"description"`
}

type GetEvsArgs struct {
	*entity.BaseCloudRequest
}
