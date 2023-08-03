package huawei

import (
	"encoding/json"
	"fmt"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/ces"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	"goutils/pkg/snowflake"
	"goutils/pkg/string_utils"
	"k8s.io/klog/v2"

	"github.com/tidwall/gjson"
)

/*
 * @author: yeshibo
 * @date: Tuesday, 2022/09/06, 8:01:43 pm
 */

const (
	// agent metric_name
	mem_usedPercent  = "mem_usedPercent"
	disk_usedPercent = "disk_usedPercent"

	// base metric_name
	disk_util_inband = "disk_util_inband"
	mem_util         = "mem_util"
)

type CesCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetCesInstance(ack accessKey.AccessKey) *CesCloudProvider {
	return &CesCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Ces,
	}
}

// GetBatchQueryMetricData 批量查询指标数据
// 注意：一次只能查询一台主机的多个指标， 因此其中的args.Metric[].Dimensions[]中的 value 必须是一样的
func (cesProvider *CesCloudProvider) GetBatchQueryMetricData(args ces.GetMetricDataArgs) ([]ces.MetricData, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    cesProvider.AccessKeyId,
		Secret: cesProvider.Secret,
	}, args.RegionId)
	projectId, err := huaweiClient.GetProjectId()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	// klog.Info(fmt.Sprintf("GetBatchQueryMetricData Request: %s", string(b)))
	cloudJson, err := huaweiClient.
		BuildUrl(serviceType.Ces, "/V1.0/{project_id}/batch-query-metric-data", args, projectId).
		Post().
		Body(b).
		SendRequest()
	if err != nil {
		return nil, err
	}
	metrics := gjson.GetBytes(cloudJson, "metrics").Array()
	if len(metrics) == 0 {
		// 调试日志
		jsonArgs, _ := json.Marshal(args)
		klog.Info(fmt.Sprintf("ack [%s] ecs [%v] GetMetricData [%v]", cesProvider.AccessKeyId, string(jsonArgs), string(cloudJson)))
	}
	// 特殊定制
	metricData := []ces.MetricData{}
	// disk_usedPercentValue := []float64{}
	// instanceId := ""
	for _, metric := range metrics {
		datapoints := metric.Get("datapoints").Array()
		for _, datapoint := range datapoints {
			metric_name := metric.Get("metric_name").String()
			// 特殊定制 将mem_usedPercent转换为mem_util，将agent的指标转换为base基础监控的指标
			// if metric_name == mem_usedPercent {
			// 	metric_name = mem_util
			// }
			// 查询包含disk_usedPercent的数据，将disk_usedPercent转换为disk_util_inband
			// 华为云缺陷： 华为云的插件获取的数据是每个盘的挂载点的使用率，而不是整个主机磁盘的使用率，且无法区分挂载点对应的盘
			// if strings.Contains(metric_name, disk_usedPercent) {
			// 	disk_usedPercentValue = append(disk_usedPercentValue, string_utils.Decimal(datapoint.Get("average").Float()))
			// 	instanceId = metric.Get("dimensions").Array()[0].Get("value").String()
			// }
			data := ces.MetricData{
				Id:         snowflake.GetWorkId(),
				InstanceId: metric.Get("dimensions").Array()[0].Get("value").String(),
				MetricName: metric_name,
				Value:      string_utils.Decimal(datapoint.Get("average").Float()),
				Timestamp:  datapoint.Get("timestamp").Int(),
				Unit:       metric.Get("unit").String(),
			}
			metricData = append(metricData, data)
		}
	}
	// 当存在disk_usedPercent时，计算disk_util_inband的最大值
	// 为啥使用最大值？因为 不管是系统盘还是数据盘，只要有一个盘的使用率超过了阈值，就需要给客户以警示作用
	// if len(disk_usedPercentValue) > 0 {
	// 	// 计算disk_usedPercent的最大值
	// 	average := string_utils.Decimal(string_utils.MaxFloat64(disk_usedPercentValue))
	// 	metricData = append(metricData, ces.MetricData{
	// 		Id:         snowflake.GetWorkId(),
	// 		InstanceId: instanceId,
	// 		MetricName: disk_util_inband,
	// 		Value:      average,
	// 		Timestamp:  time.Now().UTC().UnixMilli(),
	// 		Unit:       "%",
	// 	})
	// }
	return metricData, nil
}

func (cesProvider *CesCloudProvider) GetEcsMetric(args ces.GetMetricArgs) ([]ces.Metric, error) {
	huaweiClient := cloud_client.NewHuaweiClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    cesProvider.AccessKeyId,
		Secret: cesProvider.Secret,
	}, args.RegionId)
	projectId, err := huaweiClient.GetProjectId()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	cloudJson, err := huaweiClient.
		BuildUrl(serviceType.Ces, "/V1.0/{project_id}/metrics", args, projectId).
		Get().
		Body(b).
		SendRequest()
	if err != nil {
		return nil, err
	}
	metrics := gjson.GetBytes(cloudJson, "metrics").Array()
	if len(metrics) == 0 {
		// 调试日志
		jsonArgs, _ := json.Marshal(args)
		klog.Info(fmt.Sprintf("ack [%s] ecsMetric [%v] GetMetric [%v]", cesProvider.AccessKeyId, string(jsonArgs), string(cloudJson)))
	}
	// 特殊定制
	metricList := []ces.Metric{}
	for _, metric := range metrics {
		data := ces.Metric{
			Namespace: metric.Get("namespace").String(),
			Dimensions: []ces.Dimension{
				{
					Name:  metric.Get("dimensions").Array()[0].Get("name").String(),
					Value: metric.Get("dimensions").Array()[0].Get("value").String(),
				},
			},
			MetricName: metric.Get("metric_name").String(),
		}
		metricList = append(metricList, data)
	}
	return metricList, nil
}
