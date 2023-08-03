package huawei

import (
	"context"
	"errors"
	"fmt"
	"goutils/pkg/cloud/cloud_client"
	"goutils/pkg/cloud/cloud_request"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/entity"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/ecs"
	"goutils/pkg/cloud/entity/evs"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	timeUtil "goutils/pkg/time"
	"strings"

	"github.com/lstack-org/utils/pkg/stream"
	"github.com/tidwall/gjson"
)

const (
	CCE_NODE_TAG      = "CCE-Dynamic-Provisioning-Node"
	RESOURCE_TYPE_ECS = "Ecs"
	RESOURCE_TYPE_CCE = "Cce"
)

/*
 * @author: yeshibo
 * @date: Tuesday, 2022/09/06, 7:49:31 pm
 */

type EcsCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetEcsInstance(ack accessKey.AccessKey) *EcsCloudProvider {
	return &EcsCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Ecs,
	}
}

func (ecsProvider *EcsCloudProvider) GetEcsRequest(args ecs.GetEcsArgs) (int64, []ecs.Ecs, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    ecsProvider.AccessKeyId,
		Secret: ecsProvider.Secret,
	}, args.RegionId)
	projectId, err := huaweiClient.GetProjectId()
	if err != nil {
		return 0, nil, err
	}
	cloudJson, err := huaweiClient.
		BuildUrl(serviceType.Ecs, "/v1/{project_id}/cloudservers/detail", args, projectId).
		Get().
		SendRequest()
	if err != nil {
		return 0, nil, err
	}
	klog.V(5).Info(fmt.Sprintf("[%s] get /v1/{project_id}/cloudservers/detail : %v", ecsProvider.AccessKeyId, string(cloudJson)))
	// todo 缺少错误处理
	if gjson.GetBytes(cloudJson, "gjson.GetBytes").String() != "" {
		klog.Info("cloud response: ", string(cloudJson))
		return 0, nil, errors.New(gjson.GetBytes(cloudJson, "gjson.GetBytes").String())
	}
	evsTotal, evsList, err := GetEvsInstance(ecsProvider.AccessKey).GetEvsRequest(evs.GetEvsArgs{
		BaseCloudRequest: &entity.BaseCloudRequest{
			Context:  context.Background(),
			RegionId: args.RegionId,
		},
	})
	if err != nil {
		return 0, nil, err
	}
	if int(evsTotal) > len(evsList) {
		step := (int(evsTotal) - len(evsList)) / 1000
		for i := 2; i < (step + 1); i++ {
			_, lastEvsList, err := GetEvsInstance(ecsProvider.AccessKey).GetEvsRequest(evs.GetEvsArgs{
				BaseCloudRequest: &entity.BaseCloudRequest{
					Context:  context.Background(),
					RegionId: args.RegionId,
					Page:     i,
					PageSize: 1000,
				},
			})
			if err != nil {
				klog.Error(fmt.Sprintf("get step[%v] evs failed: %v", i, err))
				continue
			}
			evsList = append(evsList, lastEvsList...)
		}
	}

	m := stream.New(evsList).Group(func(v interface{}) interface{} { return v.(evs.Disk).InstanceId })
	var ecsList []ecs.Ecs
	count := gjson.GetBytes(cloudJson, "count").Int()
	servers := gjson.GetBytes(cloudJson, "servers").Array()
	for _, v := range servers {
		// 获取云盘信息
		systemDisk := ecs.SystemDisk{}
		var dataDisk []ecs.DataDisk
		disks := m[v.Get("id").String()]
		for _, disk := range disks {
			d := disk.(evs.Disk)
			if d.DiskType == "system" {
				systemDisk = ecs.SystemDisk{
					Size:        d.Size,
					Category:    d.Category,
					DiskName:    d.DiskName,
					DiskId:      d.Id,
					Description: d.Description,
				}
			} else {
				dataDisk = append(dataDisk, ecs.DataDisk{
					Size:        d.Size,
					Category:    d.Category,
					DiskName:    d.DiskName,
					DiskId:      d.Id,
					Description: d.Description,
				})
			}
		}
		// 安全组
		var securityGroups []string
		for _, r := range v.Get("security_groups").Array() {
			securityGroups = append(securityGroups, r.Get("id").String())
		}
		// ip
		privateIp := ""
		eipIps := []string{}
		for _, addresses := range v.Get("addresses").Get(v.Get("metadata.vpc_id").String()).Array() {
			if addresses.Get("OS-EXT-IPS:type").String() == "fixed" {
				privateIp = addresses.Get("addr").String()
			} else {
				eipIps = append(eipIps, addresses.Get("addr").String())
			}
		}
		t := timeUtil.StringToGoTime(v.Get("created").String())
		updated := timeUtil.StringToGoTime(v.Get("updated").String())
		regionName := cloud_request.GetRegionName(args.RegionId)
		resourceType := RESOURCE_TYPE_ECS
		for _, tag := range v.Get("tags").Array() {
			if find := strings.Contains(tag.String(), CCE_NODE_TAG); find {
				resourceType = RESOURCE_TYPE_CCE
			}
		}
		e := ecs.Ecs{
			Id:                v.Get("id").String(),
			Name:              v.Get("name").String(),
			InstanceName:      v.Get("name").String(),
			RegionId:          args.RegionId,
			CloudProviderId:   "huaweiyun",
			CloudProviderName: "HCSO",
			RegionName:        regionName,
			AckUid:            args.AckUid,
			Status:            v.Get("status").String(),
			OsType:            v.Get("metadata.os_type").String(),
			InstanceType:      v.Get("flavor.id").String(),
			Description:       v.Get("description").String(),
			CPU:               v.Get("flavor.vcpus").Int(),
			Memory:            v.Get("flavor.ram").Int() / 1024,
			SystemDisk:        systemDisk,
			DataDisks:         dataDisk,
			Vpc: ecs.Vpc{
				VpcId: v.Get("metadata.vpc_id").String(),
			},
			SecurityGroups: securityGroups,
			CreateTime:     t,
			UpdateTime:     updated,
			EipIp:          eipIps,
			PrivateIp:      privateIp,
			ResourceType:   resourceType,
		}
		ecsList = append(ecsList, e)
	}
	return count, ecsList, nil
}
