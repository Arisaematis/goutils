package entity

import (
	"context"
)

type BaseCloudRequest struct {
	context.Context `json:"-"`
	RegionId        string `json:"-"`
	AckUid          string `json:"-"`
	Page            int    `json:"-" huawei:"offset,omitempty"`
	PageSize        int    `json:"-" huawei:"limit,omitempty"`
	ServiceId       string `json:"-" huawei:"service_id,omitempty"`
}

type BaseCloudResponse struct {
	Id           string `json:"id" bson:"id"`
	Name         string `json:"name" bson:"name"`
	ResourceType string `json:"resourceType" bson:"resourceType"`
	Description  string `json:"description" bson:"description"`
	RegionId     string `json:"regionId" bson:"regionId"`
	AckUid       string `json:"ackUid" bson:"ackUid"`
	Status       string `json:"status" bson:"status"`
}

// Resource 资源
type Resource struct {
	BaseCloudResponse `json:"inline" bson:"inline"`
	Data              interface{} `json:"data" bson:"data"`
}
