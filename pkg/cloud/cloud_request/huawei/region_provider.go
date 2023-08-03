package huawei

import (
	"fmt"
	"goutils/pkg/cloud/cloud_client"
	"goutils/pkg/cloud/cloud_request"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/region"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	"goutils/pkg/setting"

	"github.com/tidwall/gjson"
)

type RegionCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetRegionInstance(ack accessKey.AccessKey) *RegionCloudProvider {
	return &RegionCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Region,
	}
}

func (regionProvider *RegionCloudProvider) GetRegionRequest(args entity.BaseCloudRequest) ([]region.Region, error) {
	// 判断 redis 中是否存在
	baseUrl := ""
	if setting.EnvConfig.ClientType == "huawei" {
		baseUrl = "https://iam.myhuaweicloud.com/v3/regions"
	} else {
		if regionProvider.AccessKeyType == 0 {
			baseUrl = cloud_request.HCSO_IAM_ENDPOINT + "/v3/regions"
		} else {
			baseUrl = cloud_request.HCSO_OUT_IAM_ENDPOINT + "/v3/regions"
		}
	}
	cloudJson, err := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    regionProvider.AccessKeyId,
		Secret: regionProvider.Secret,
	}, args.RegionId).Get().Url(baseUrl).SendRequest()
	if err != nil {
		return nil, err
	}
	klog.Info(fmt.Sprintf("request url [%s] response[%v]", baseUrl, string(cloudJson)))
	var regions []region.Region
	gRegions := gjson.GetBytes(cloudJson, "regions").Array()
	for _, v := range gRegions {
		regionId := v.Get("id").String()
		regionName := v.Get("locales.zh-cn").String()
		r := region.Region{
			RegionId:   regionId,
			RegionName: regionName,
		}
		regions = append(regions, r)
	}
	return regions, nil
}

func (regionProvider *RegionCloudProvider) GetEndPointRequest(args entity.BaseCloudRequest) ([]interface{}, error) {
	baseUrl := ""
	if setting.EnvConfig.ClientType == "huawei" {
		baseUrl = "https://iam.myhuaweicloud.com/v3/endpoints"
	} else {
		if regionProvider.AccessKeyType == 0 {
			baseUrl = cloud_request.HCSO_IAM_ENDPOINT + "/v3/endpoints"
		} else {
			baseUrl = cloud_request.HCSO_OUT_IAM_ENDPOINT + "/v3/endpoints"
		}
	}
	_, err := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    regionProvider.AccessKeyId,
		Secret: regionProvider.Secret,
	}, args.RegionId).Get().Url(baseUrl).SendRequest()
	return nil, err
}
