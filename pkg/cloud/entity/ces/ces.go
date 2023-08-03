package ces

import (
	"goutils/pkg/cloud/entity"
)

type MetricData struct {
	Id         string  `bson:"id"`
	InstanceId string  `bson:"instance_id"`
	MetricName string  `bson:"metric_name"`
	Value      float64 `bson:"value"`
	Timestamp  int64   `bson:"timestamp"`
	Unit       string  `bson:"unit"`
}

type GetMetricDataArgs struct {
	*entity.BaseCloudRequest
	Metrics []Metric `json:"metrics"`
	From    int64    `json:"from"`
	To      int64    `json:"to"`
	Period  string   `json:"period"`
	Filter  string   `json:"filter"`
}

type GetMetricArgs struct {
	*entity.BaseCloudRequest
	Dim0 string `huawei:"dim.0"`
}

type Metric struct {
	Namespace  string      `json:"namespace"`
	Dimensions []Dimension `json:"dimensions"`
	MetricName string      `json:"metric_name"`
}

type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
