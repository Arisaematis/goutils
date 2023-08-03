package ecs

import (
	"goutils/pkg/cloud/entity"
	"time"
)

type Ecs struct {
	Id                string     `json:"id" bson:"id"`
	Name              string     `json:"name,omitempty" bson:"name,omitempty"`
	InstanceName      string     `json:"instanceName" bson:"instanceName"`
	RegionId          string     `json:"regionId" bson:"regionId"`
	CloudProviderId   string     `json:"cloudProviderId" bson:"cloudProviderId"`
	CloudProviderName string     `json:"cloudProviderName" bson:"cloudProviderName"`
	RegionName        string     `json:"regionName" bson:"regionName"`
	AckUid            string     `json:"ackUid,omitempty" bson:"ackUid,omitempty"`
	Status            string     `json:"status" bson:"status"`
	OsType            string     `json:"osType" bson:"osType"`
	InstanceType      string     `json:"instanceType" bson:"instanceType"`
	Description       string     `json:"description" bson:"description"`
	CPU               int64      `json:"cpu" bson:"cpu"`
	Memory            int64      `json:"memory" bson:"memory"`
	SystemDisk        SystemDisk `json:"systemDisk" bson:"systemDisk"`
	DataDisks         []DataDisk `json:"dataDisks" bson:"dataDisks"`
	PublicIp          string     `json:"publicIp" bson:"publicIp"`
	PrivateIp         string     `json:"privateIp" bson:"privateIp"`
	EipIp             []string   `json:"eipIp" bson:"eipIp"`
	Eip               []Eip      `json:"eip" bson:"eip"`
	Vpc               Vpc        `json:"vpc" bson:"vpc"`
	SecurityGroups    []string   `json:"securityGroups" bson:"securityGroups"`
	CreateTime        time.Time  `json:"createTime" bson:"createTime,omitempty"`
	UpdateTime        time.Time  `json:"updateTime" bson:"updateTime,omitempty"`
	ResourceType      string     `json:"resourceType" bson:"resourceType"`
	// Tags              []Tags     `json:"tags"`
}

type SystemDisk struct {
	Size        int    `json:"size" bson:"size"`
	Category    string `json:"category" bson:"category"`
	DiskName    string `json:"diskName" bson:"diskName"`
	DiskId      string `json:"diskId" bson:"diskId"`
	Description string `json:"description" bson:"description"`
}
type DataDisk struct {
	Size        int    `json:"size" bson:"size"`
	Category    string `json:"category" bson:"category"`
	DiskName    string `json:"diskName" bson:"diskName"`
	DiskId      string `json:"diskId" bson:"diskId"`
	Description string `json:"description" bson:"description"`
}
type Eip struct {
	EipId    string `json:"eipId" bson:"eipId"`
	EipName  string `json:"eipName" bson:"eipName"`
	EipIp    string `json:"eipIp" bson:"eipIp"`
	BandWith int    `json:"bandWith" bson:"bandWith"`
}
type Vpc struct {
	VpcId     string `json:"vpcId" bson:"vpcId"`
	VswitchId string `json:"vswitchId" bson:"vswitchId"`
}
type Tags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetEcsArgs struct {
	*entity.BaseCloudRequest
	Page     int `json:"-" huawei:"offset"`
	PageSize int `json:"-" huawei:"limit"`
}
