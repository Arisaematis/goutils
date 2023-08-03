package cloud_request

import (
	"github.com/patrickmn/go-cache"
	"goutils/pkg/cloud/cloud_request/service_type"
)

const (
	HCSO_IAM_ENDPOINT = "https://iam-pub.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_ECS_ENDPOINT = "https://ecs.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_CES_ENDPOINT = "https://ces.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_EVS_ENDPOINT = "https://evs.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_VPC_ENDPOINT = "https://vpc.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_IMS_ENDPOINT = "https://ims.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_OBS_ENDPOINT = "https://obs.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_SFS_ENDPOINT = "https://sfs-turbo.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_RDS_ENDPOINT = "https://rds.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"
	HCSO_EIP_ENDPOINT = "https://eip.cn-north-6046.taiyuanxc04.gov.huaweicloud.com"

	HCSO_OUT_IAM_ENDPOINT = "https://iam-pub.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_ECS_ENDPOINT = "https://ecs.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_CES_ENDPOINT = "https://ces.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_EVS_ENDPOINT = "https://evs.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_VPC_ENDPOINT = "https://vpc.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_IMS_ENDPOINT = "https://ims.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_OBS_ENDPOINT = "https://obs.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_SFS_ENDPOINT = "https://sfs-turbo.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_RDS_ENDPOINT = "https://rds.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"
	HCSO_OUT_EIP_ENDPOINT = "https://eip.cn-north-6045.taiyuanxc03.gov.huaweicloud.com"

	HCSO_REGION_TAIYUAN        = "政务外网"
	HCSO_REGION_TAIYUAN_ID     = "cn-north-6046"
	HCSO_OUT_REGION_TAIYUAN    = "互联网"
	HCSO_OUT_REGION_TAIYUAN_ID = "cn-north-6045"
)

var HCSO_ENDPOINT_MAP map[service_type.Type]string
var REGION_MAP map[string]string
var DefaultNoExpirationCache *cache.Cache
var HCSO_ENDPOINT_MAP_FOR_REGION map[string]map[service_type.Type]string

func init() {
	DefaultNoExpirationCache = cache.New(cache.DefaultExpiration, 0)
	HCSO_ENDPOINT_MAP = map[service_type.Type]string{}

	REGION_MAP = make(map[string]string)
	REGION_MAP[HCSO_REGION_TAIYUAN_ID] = HCSO_REGION_TAIYUAN
	REGION_MAP[HCSO_OUT_REGION_TAIYUAN_ID] = HCSO_OUT_REGION_TAIYUAN

	HCSO_ENDPOINT_MAP_NEW1 := map[service_type.Type]string{}
	HCSO_ENDPOINT_MAP_FOR_REGION = map[string]map[service_type.Type]string{}
	HCSO_ENDPOINT_MAP_NEW1[service_type.Iam] = HCSO_OUT_IAM_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Ecs] = HCSO_OUT_ECS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Region] = HCSO_OUT_IAM_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Ces] = HCSO_OUT_CES_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Evs] = HCSO_OUT_EVS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Vpc] = HCSO_OUT_VPC_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Obs] = HCSO_OUT_OBS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Sfs] = HCSO_OUT_SFS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Rds] = HCSO_OUT_RDS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW1[service_type.Ims] = HCSO_OUT_IMS_ENDPOINT
	HCSO_ENDPOINT_MAP_FOR_REGION[HCSO_OUT_REGION_TAIYUAN_ID] = HCSO_ENDPOINT_MAP_NEW1
	HCSO_ENDPOINT_MAP_NEW2 := map[service_type.Type]string{}
	HCSO_ENDPOINT_MAP_NEW2[service_type.Iam] = HCSO_IAM_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Ecs] = HCSO_ECS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Region] = HCSO_IAM_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Ces] = HCSO_CES_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Evs] = HCSO_EVS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Vpc] = HCSO_VPC_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Obs] = HCSO_OBS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Sfs] = HCSO_SFS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Rds] = HCSO_RDS_ENDPOINT
	HCSO_ENDPOINT_MAP_NEW2[service_type.Ims] = HCSO_IMS_ENDPOINT
	HCSO_ENDPOINT_MAP_FOR_REGION[HCSO_REGION_TAIYUAN_ID] = HCSO_ENDPOINT_MAP_NEW2
}

func GetRegionName(regionId string) string {
	s, ok := REGION_MAP[regionId]
	regionName := ""
	if ok {
		regionName = s
	}
	return regionName
}
