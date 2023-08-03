package huawei

import (
	"fmt"
	"goutils/pkg/cloud/cloud_client"
	serviceType "goutils/pkg/cloud/cloud_request/service_type"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"goutils/pkg/cloud/entity/obs"
	"goutils/pkg/cloud/entity/request"
	huaweiSign "goutils/pkg/cloud/sdk_sign/huawei"
	"k8s.io/klog/v2"
	"strconv"

	huaweiObs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

const (
	HeaderXObsBucketType      = "x-obs-bucket-type"
	HeaderXObsBucketTypeValue = "OBJECT"
	HeaderXObsDate            = "x-obs-date"
	ObsDateFormat             = "Mon, _2 Jan 2006 15:04:05 MST"
)

type ObsCloudProvider struct {
	accessKey.AccessKey
	serviceType.Type
}

func GetObsInstance(ack accessKey.AccessKey) *ObsCloudProvider {
	return &ObsCloudProvider{
		AccessKey: ack,
		Type:      serviceType.Obs,
	}
}

func (obsProvider *ObsCloudProvider) GetObsRequest(args obs.GetObsArgs) (int, []obs.Bucket, error) {
	huaweiObsClient, err := cloud_client.NewHuaweiObsClient(args.Context, &huaweiSign.HuaweiSigner{
		Key:    obsProvider.AccessKeyId,
		Secret: obsProvider.Secret,
	}, args.RegionId)
	if err != nil {
		return 0, nil, err
	}
	huaweiBuckets, err := huaweiObsClient.ListBuckets(nil)
	if err != nil {
		return 0, nil, err
	}
	buckets := []obs.Bucket{}
	for _, huaweiBucket := range huaweiBuckets.Buckets {
		bucket := obs.Bucket{
			BaseCloudResponse: request.BaseCloudResponse{
				Id:                fmt.Sprintf("%s-%s", huaweiBucket.Name, obsProvider.AckUid),
				Name:              huaweiBucket.Name,
				AckUid:            obsProvider.AckUid,
				ResourceType:      "Obs",
				CloudProviderId:   obsProvider.CloudProviderId,
				CloudProviderName: "华为HCSO",
				CreateTime:        huaweiBucket.CreationDate,
			},
			BucketType: huaweiBucket.BucketType,
		}
		// 查询bucket 存储量
		storageSize, err := huaweiObsClient.GetBucketStorageInfo(huaweiBucket.Name)
		if err != nil {
			klog.Error(fmt.Sprintf("get [%s] BucketStorageInfo err : %v", huaweiBucket.Name, err))
		} else {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(storageSize.Size)/float64(1024)/float64(1024)), 64)
			bucket.Size = value
		}
		// 查询bucket 限制额度
		quota, err := huaweiObsClient.GetBucketQuota(huaweiBucket.Name)
		if err != nil {
			klog.Error(fmt.Sprintf("get [%s] BucketQuota err : %v", huaweiBucket.Name, err))
		} else {
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(quota.Quota)/float64(1024)/float64(1024)), 64)
			bucket.QuotaSize = value
		}
		buckets = append(buckets, bucket)
	}
	return len(buckets), buckets, nil
}

func ListObjects(huaweiObsClient *huaweiObs.ObsClient, bucketName, marker string) (int64, error) {
	// 查询bucket 下的Object文件量
	// 文件对象默认返回 1000个
	huaweiObject, err := huaweiObsClient.ListObjects(&huaweiObs.ListObjectsInput{
		Bucket: bucketName,
		Marker: marker,
	})
	if err != nil {
		klog.Error(fmt.Sprintf("ListObjects [%s] err: %v", bucketName, err))
		return 0, err
	}
	var objectSize int64
	for _, content := range huaweiObject.Contents {
		objectSize += content.Size
	}
	// 判断是否列举完成
	if huaweiObject.IsTruncated {
		i, err := ListObjects(huaweiObsClient, bucketName, huaweiObject.NextMarker)
		if err != nil {
			klog.Error(fmt.Sprintf("recursion ListObjects [%s] err: %v", bucketName, err))
			return 0, err
		}
		objectSize += i
	}
	return objectSize, nil
}
