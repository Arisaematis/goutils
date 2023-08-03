package huawei

import (
	"github.com/tidwall/gjson"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/ims"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
)

type ImsCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetImsInstance(ack accessKey.AccessKey) *ImsCloudProvider {
	return &ImsCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Ims,
	}
}

func (imsProvider *ImsCloudProvider) GetImsRequest(args ims.GetImsRequest) (int64, []ims.Ims, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    imsProvider.AccessKeyId,
		Secret: imsProvider.Secret,
	}, args.RegionId)
	cloudJson, err := huaweiClient.
		BuildUrlWithProject(imsProvider.Type, "/v2/cloudimages", args).
		Get().
		// SetResponseErrorHandler(ResponseErrorHandler).
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	imsList := []ims.Ims{}
	images := gjson.GetBytes(cloudJson, "images").Array()
	for _, v := range images {
		imsList = append(imsList, ims.Ims{
			Resource: entity.Resource{
				BaseCloudResponse: entity.BaseCloudResponse{
					Id:           v.Get("id").String(),
					Name:         v.Get("name").String(),
					ResourceType: serviceType.Ims.String(),
					Description:  v.Get("description").String(),
					RegionId:     args.RegionId,
					AckUid:       imsProvider.AckUid,
					Status:       v.Get("status").String(),
				},
				Data: ims.ImsData{},
			},
		})
	}
	return int64(len(imsList)), imsList, nil
}
