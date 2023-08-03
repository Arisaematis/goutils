package huawei

import (
	"github.com/tidwall/gjson"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/vpc"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
)

type VpcCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetVpcInstance(ack accessKey.AccessKey) *VpcCloudProvider {
	return &VpcCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Vpc,
	}
}

func (vpcProvider *VpcCloudProvider) GetVpcRequest(args vpc.GetVpcRequest) (int64, []vpc.Vpc, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    vpcProvider.AccessKeyId,
		Secret: vpcProvider.Secret,
	}, args.RegionId)
	cloudJson, err := huaweiClient.
		BuildUrlWithProject(vpcProvider.Type, "/v3/{project_id}/vpc/vpcs", args).
		Get().
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	vpcList := []vpc.Vpc{}
	// todo vpc 数据可能不准确，华为云vpc 默认拉取2000条
	vpcs := gjson.GetBytes(cloudJson, "vpcs").Array()
	for _, v := range vpcs {
		vpcList = append(vpcList, vpc.Vpc{
			Resource: entity.Resource{
				BaseCloudResponse: entity.BaseCloudResponse{
					Id:           v.Get("id").String(),
					Name:         v.Get("name").String(),
					ResourceType: serviceType.Vpc.String(),
					Description:  v.Get("description").String(),
					RegionId:     args.RegionId,
					AckUid:       vpcProvider.AckUid,
				},
				Data: vpc.VpcData{
					CidrBlock: v.Get("cidr").String(),
				},
			},
		})
	}
	return int64(len(vpcList)), vpcList, nil
}
