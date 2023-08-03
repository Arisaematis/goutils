package huawei

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/evs"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	"k8s.io/klog/v2"
)

type EvsCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetEvsInstance(ack accessKey.AccessKey) *EvsCloudProvider {
	return &EvsCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Ces,
	}
}

func (provider *EvsCloudProvider) GetEvsRequest(args evs.GetEvsArgs) (int64, []evs.Disk, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    provider.AccessKeyId,
		Secret: provider.Secret,
	}, args.RegionId)
	projectId, err := huaweiClient.GetProjectId()
	if err != nil {
		return 0, nil, err
	}
	cloudJson, err := huaweiClient.
		BuildUrl(serviceType.Evs, "/v2/{project_id}/cloudvolumes/detail", args, projectId).
		Get().
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	klog.V(5).Info(fmt.Sprintf("[%s] get [/v2/{project_id}/cloudvolumes/detail] : %v", provider.AccessKeyId, string(cloudJson)))
	// 异常处理
	if gjson.GetBytes(cloudJson, "gjson.GetBytes").String() != "" {
		klog.Info("cloud response: ", string(cloudJson))
		return 0, nil, errors.New(gjson.GetBytes(cloudJson, "gjson.GetBytes").String())
	}
	diskList := []evs.Disk{}
	count := gjson.GetBytes(cloudJson, "count").Int()
	volumes := gjson.GetBytes(cloudJson, "volumes").Array()
	for _, v := range volumes {
		s := v.Get("status").String()
		r := v.Get("attachments").Array()
		instanceId := ""
		if s == "in-use" && len(r) > 0 {
			instanceId = r[0].Get("server_id").String()
		}
		// 公有云上也可以如此，私有云上不行
		// bootable := v.Get("bootable").String()
		// diskType := "data"
		// if bootable == "true" {
		// 	diskType = "system"
		// }
		// 私有云根据volume_image_metadata.__os_version判断
		diskType := "data"
		volume_image_metadata := v.Get("volume_image_metadata")
		if volume_image_metadata.Exists() {
			osVersion := volume_image_metadata.Get("__os_version")
			if osVersion.Exists() {
				diskType = "system"
			}
		}
		disk := evs.Disk{
			Id:          v.Get("id").String(),
			Size:        int(v.Get("size").Int()),
			Category:    v.Get("volume_type").String(),
			DiskName:    v.Get("name").String(),
			InstanceId:  instanceId,
			DiskType:    diskType,
			Description: v.Get("description").String(),
		}
		diskList = append(diskList, disk)
	}
	return count, diskList, nil
}
