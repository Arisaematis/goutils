package cloud_request

import (
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/ces"
	"goutils/pkg/cloud/entity/ecs"
	"goutils/pkg/cloud/entity/eip"
	"goutils/pkg/cloud/entity/evs"
	"goutils/pkg/cloud/entity/ims"
	"goutils/pkg/cloud/entity/obs"
	"goutils/pkg/cloud/entity/rds"
	"goutils/pkg/cloud/entity/region"
	"goutils/pkg/cloud/entity/sfs"
	"goutils/pkg/cloud/entity/vpc"
)

/*
 * @author: yeshibo
 * @date: Tuesday, 2022/09/06, 10:24:38 am
 */

// IRequestProducerFactory
// 云服务资源提供者的生产工厂
type IRequestProducerFactory interface {
	GetEcsRequestProducer(accessKey.AccessKey) BaseEcsRequestProducer
	GetRegionRequestProducer(accessKey.AccessKey) BaseRegionRequestProducer
	GetCESRequestProducer(accessKey.AccessKey) BaseCesRequestProducer
	GetEvsRequestProducer(accessKey.AccessKey) BaseEvsRequestProducer
	GetObsRequestProducer(accessKey.AccessKey) BaseObsRequestProducer
	GetNasRequestProducer(accessKey.AccessKey) BaseNasRequestProducer
	GetEipRequestProducer(accessKey.AccessKey) BaseEipRequestProducer
	GetImsRequestProducer(accessKey.AccessKey) BaseImsRequestProducer
	GetVpcRequestProducer(accessKey.AccessKey) BaseVpcRequestProducer
	GetRdsRequestProducer(accessKey.AccessKey) BaseRdsRequestProducer
}

// BaseRegionRequestProducer 地域查询提供者
type BaseRegionRequestProducer interface {
	GetRegionRequest(entity.BaseCloudRequest) ([]region.Region, error)
	GetEndPointRequest(entity.BaseCloudRequest) ([]interface{}, error)
}

// BaseCesRequestProducer 云监控
type BaseCesRequestProducer interface {
	// 此接口用于定时任务同步监控数据
	GetBatchQueryMetricData(ces.GetMetricDataArgs) ([]ces.MetricData, error)
	// 查询指标列表
	GetEcsMetric(ces.GetMetricArgs) ([]ces.Metric, error)
}

// BaseEcsRequestProducer 云主机
type BaseEcsRequestProducer interface {
	GetEcsRequest(ecs.GetEcsArgs) (int64, []ecs.Ecs, error)
}

type BaseEvsRequestProducer interface {
	GetEvsRequest(evs.GetEvsArgs) (int64, []evs.Disk, error)
}

type BaseObsRequestProducer interface {
	GetObsRequest(obs.GetObsArgs) (int, []obs.Bucket, error)
}

type BaseNasRequestProducer interface {
	GetNasRequest(sfs.GetNasRequest) (int, []sfs.Nas, error)
}

type BaseEipRequestProducer interface {
	GetEipRequest(eip.GetEipRequest) (int64, []eip.Eip, error)
}

type BaseImsRequestProducer interface {
	GetImsRequest(ims.GetImsRequest) (int64, []ims.Ims, error)
}

type BaseVpcRequestProducer interface {
	GetVpcRequest(vpc.GetVpcRequest) (int64, []vpc.Vpc, error)
}

type BaseRdsRequestProducer interface {
	GetRdsRequest(rds.GetRdsRequest) (int64, []rds.Rds, error)
}
