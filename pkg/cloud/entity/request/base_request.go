package request

import (
	"goutils/pkg/cloud/entity/accesskey"
	"time"
)

type BaseCloudResponse struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	AckUid            string    `json:"ackUid"`
	ResourceType      string    `json:"resourceType"`
	CloudProviderId   string    `json:"cloudProviderId"`
	CloudProviderName string    `json:"cloudProviderName"`
	RegionId          string    `json:"regionId"`
	RegionName        string    `json:"regionName"`
	CreateTime        time.Time `json:"crateTime"`
}

func NewBaseCloudResponse(id, name, resourceType, regionId string, createTime time.Time, accessKey accesskey.AccessKey) BaseCloudResponse {
	return BaseCloudResponse{
		Id:                id,
		Name:              name,
		AckUid:            accessKey.AckUid,
		ResourceType:      resourceType,
		CloudProviderId:   accessKey.CloudProviderId,
		CloudProviderName: "",
		RegionId:          regionId,
		RegionName:        "",
		CreateTime:        createTime,
	}
}
