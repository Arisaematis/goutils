package huawei

import (
	"github.com/tidwall/gjson"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/eip"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
)

type EipCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetEipInstance(ack accessKey.AccessKey) *EipCloudProvider {
	return &EipCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Eip,
	}
}

func (eipProvider *EipCloudProvider) GetEipRequest(args eip.GetEipRequest) (int64, []eip.Eip, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    eipProvider.AccessKeyId,
		Secret: eipProvider.Secret,
	}, args.RegionId)
	cloudJson, err := huaweiClient.
		BuildUrlWithProject(serviceType.Vpc, "/v1/{project_id}/publicips", args).
		Get().
		// SetResponseErrorHandler(ResponseErrorHandler).
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	var eipList []eip.Eip
	publicips := gjson.GetBytes(cloudJson, "publicips").Array()
	for _, v := range publicips {
		eipList = append(eipList, eip.Eip{
			Resource: entity.Resource{
				BaseCloudResponse: entity.BaseCloudResponse{
					Id:           v.Get("id").String(),
					Name:         v.Get("name").String(),
					ResourceType: serviceType.Eip.String(),
					Description:  v.Get("description").String(),
					RegionId:     args.RegionId,
					AckUid:       eipProvider.AckUid,
					Status:       v.Get("status").String(),
				},
				Data: eip.EipData{},
			},
		})
	}
	return int64(len(eipList)), eipList, nil
}
