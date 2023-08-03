package huawei

import (
	"encoding/json"
	"fmt"
	"goutils/pkg/cloud/cloud_client"
	"goutils/pkg/cloud/cloud_request"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/request"
	"goutils/pkg/cloud/entity/sfs"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	timeUtil "goutils/pkg/time"
	"k8s.io/klog/v2"
	"math"
	"strconv"

	"github.com/tidwall/gjson"
)

type NasCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetNasInstance(ack accessKey.AccessKey) *NasCloudProvider {
	return &NasCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Obs,
	}
}

func (nasProvider *NasCloudProvider) GetNasRequest(args sfs.GetNasRequest) (int, []sfs.Nas, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    nasProvider.AccessKeyId,
		Secret: nasProvider.Secret,
	}, args.RegionId)
	projectId, err := huaweiClient.GetProjectId()
	if err != nil {
		return 0, nil, err
	}
	b, err := json.Marshal(args)
	if err != nil {
		return 0, nil, err
	}
	// klog.Info(fmt.Sprintf("GetBatchQueryMetricData Request: %s", string(b)))
	cloudJson, err := huaweiClient.
		BuildUrl(serviceType.Sfs, "/v1/{project_id}/sfs-turbo/shares/detail", args, projectId).
		Get().
		Body(b).
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	shares := gjson.GetBytes(cloudJson, "shares").Array()
	count := gjson.GetBytes(cloudJson, "count").Int()
	klog.Info(fmt.Sprintf("[%s] get [%s] sfs: %v", nasProvider.AccessKeyId, args.RegionId, string(cloudJson)))
	// 特殊定制
	nasData := []sfs.Nas{}
	for _, share := range shares {
		regionId := share.Get("region").String()
		size, err := strconv.ParseFloat(share.Get("size").String(), 64)
		if err != nil {
			klog.Error(fmt.Sprintf("sfs's size [%s] ParseFloat err: %v", share.Get("size").String(), err))
		}
		avail_capacity, err := strconv.ParseFloat(share.Get("avail_capacity").String(), 64)
		if err != nil {
			klog.Error(fmt.Sprintf("sfs's avail_capacity [%s] ParseFloat err: %v", share.Get("avail_capacity").String(), err))
		}
		data := sfs.Nas{
			BaseCloudResponse: request.BaseCloudResponse{
				Id:                share.Get("id").String(),
				Name:              share.Get("name").String(),
				AckUid:            nasProvider.AckUid,
				ResourceType:      "Nas",
				CloudProviderId:   nasProvider.CloudProviderId,
				CloudProviderName: "华为HCSO",
				RegionId:          regionId,
				RegionName:        cloud_request.GetRegionName(regionId),
				CreateTime:        timeUtil.StringToGoTime(share.Get("created_at").String()),
			},
			// 将 GB 统一转换成 MB
			Size:      Decimal((size - avail_capacity) * 1024),
			QuotaSize: Decimal(size * 1024),
		}
		nasData = append(nasData, data)
	}
	return int(count), nasData, nil
}

func Decimal(num float64) float64 {
	if math.IsNaN(num) {
		return 0
	}
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
