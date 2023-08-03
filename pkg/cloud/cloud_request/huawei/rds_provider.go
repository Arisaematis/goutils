package huawei

import (
	"github.com/tidwall/gjson"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/rds"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
)

type RdsCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetRdsInstance(ack accessKey.AccessKey) *RdsCloudProvider {
	return &RdsCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Rds,
	}
}

// func ResponseErrorHandler(res []byte) error {
// 	return nil
// }

func (rdsProvider *RdsCloudProvider) GetRdsRequest(args rds.GetRdsRequest) (int64, []rds.Rds, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    rdsProvider.AccessKeyId,
		Secret: rdsProvider.Secret,
	}, args.RegionId)
	cloudJson, err := huaweiClient.
		BuildUrlWithProject(rdsProvider.Type, "/v3/{project_id}/instances", args).
		Get().
		// SetResponseErrorHandler(ResponseErrorHandler).
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	rdsList := []rds.Rds{}
	instances := gjson.GetBytes(cloudJson, "instances").Array()
	total := gjson.GetBytes(cloudJson, "total_count").Int()
	for _, v := range instances {
		rdsList = append(rdsList, rds.Rds{
			Resource: entity.Resource{
				BaseCloudResponse: entity.BaseCloudResponse{
					Id:           v.Get("id").String(),
					Name:         v.Get("name").String(),
					ResourceType: serviceType.Rds.String(),
					Description:  v.Get("description").String(),
					RegionId:     args.RegionId,
					AckUid:       rdsProvider.AckUid,
					Status:       v.Get("status").String(),
				},
				Data: rds.RdsData{},
			},
		})
	}
	return total, rdsList, nil
}
