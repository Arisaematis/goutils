package cloud_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"goutils/pkg/cloud/cloud_request"
	"goutils/pkg/cloud/cloud_request/service_type"
	"goutils/pkg/cloud/http/query"
	sdkSign "goutils/pkg/cloud/sdk_sign"
	"goutils/pkg/setting"
	"goutils/pkg/string_utils"
	"k8s.io/klog/v2"
	"net/url"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/tidwall/gjson"
)

type HuaweiClient struct {
	*BaseClient
	RegionId string
}

func NewHuaweiClient(cxt context.Context, huaweiSign sdkSign.ISign, regionId string) *HuaweiClient {
	return &HuaweiClient{
		BaseClient: NewBaseClientWithContext(cxt, huaweiSign),
		RegionId:   regionId,
	}
}

func NewHuaweiObsClient(cxt context.Context, huaweiSign sdkSign.ISign, regionId string) (*obs.ObsClient, error) {
	huaweiClient := NewHuaweiClient(cxt, huaweiSign, regionId)
	endPoint := huaweiClient.BuildEndpoint(service_type.Obs, "", nil)
	return obs.New(huaweiSign.GetSecretKey(), huaweiSign.GetSecret(), endPoint)
}

// GetProjectId 查询华为云 projectId
// todo 加上缓存 提升整体的查询速度
func (huaweiClient *HuaweiClient) GetProjectId() (string, error) {
	// 1. 从缓存中查询 全部的 projectId
	secretKey := huaweiClient.GetSecretKey()
	cacheKey := secretKey + huaweiClient.RegionId
	keyValue, found := cloud_request.DefaultNoExpirationCache.Get(cacheKey)
	if found || keyValue != nil {
		return fmt.Sprintf("%s", keyValue), nil
	}
	baseUrl := ""
	if setting.EnvConfig.ClientType == "huawei" {
		baseUrl = "https://iam.myhuaweicloud.com/v3/auth/projects"
	} else {
		// if setting.AppConfig.Network == 0 {
		// 	baseUrl = cloud_request.HCSO_IAM_ENDPOINT + "/v3/auth/projects"
		// } else {
		// 	baseUrl = cloud_request.HCSO_OUT_IAM_ENDPOINT + "/v3/auth/projects"
		// }
		regionId := huaweiClient.RegionId
		m, ok := cloud_request.HCSO_ENDPOINT_MAP_FOR_REGION[regionId]
		if !ok {
			return "", fmt.Errorf("HCSO RegionId[%s] not exist", regionId)
		}
		baseUrl = m[service_type.Iam] + "/v3/auth/projects"
	}
	var projectId string
	cloudJson, err := huaweiClient.Get().
		Url(baseUrl).
		SendRequest()
	if err != nil {
		return projectId, err
	}
	// 打印 projectId json 数据
	klog.V(6).Info(fmt.Sprintf("get project result: %v", string(cloudJson)))
	gProjects := gjson.GetBytes(cloudJson, "projects").Array()
	for _, v := range gProjects {
		name := v.Get("name").String()
		if huaweiClient.RegionId == name {
			// 存入缓存
			cloud_request.DefaultNoExpirationCache.Set(cacheKey, v.Get("id").String(), cache.DefaultExpiration)
			return v.Get("id").String(), nil
		}
	}
	return projectId, errors.New("projectId not find, Please check whether the current requested domain has access rights")
}

func (huaweiClient *HuaweiClient) BuildEndpoint(serviceType service_type.Type,
	baseUrl string, args interface{}, path ...string) string {
	huaweiEndpoint := ""
	if setting.EnvConfig.ClientType == "huawei" {
		if serviceType == service_type.Obs {
			huaweiEndpoint = fmt.Sprintf("https://%s.myhuaweicloud.com%s", serviceType.String(), BuildBaseUrl(baseUrl, args, path...))
		} else {
			// 华为 多项目（projectId） 时 需要在header 中带上 X-Project-Id
			// 添加需要签名的其他头域，或者其他用途的头域，如多项目场景中添加X-Project-Id，或者全局服务场景中添加X-Domain-Id。
			huaweiEndpoint = fmt.Sprintf("https://%s.%s.myhuaweicloud.com%s", serviceType.String(),
				strings.Split(huaweiClient.RegionId, "_")[0], BuildBaseUrl(baseUrl, args, path...))
		}
	} else {
		regionId := huaweiClient.RegionId
		m, ok := cloud_request.HCSO_ENDPOINT_MAP_FOR_REGION[regionId]
		var s string
		if !ok {
			s = cloud_request.HCSO_ENDPOINT_MAP[serviceType]
		} else {
			s = m[serviceType]
		}
		huaweiEndpoint = fmt.Sprintf("%s%s", s, BuildBaseUrl(baseUrl, args, path...))
	}
	return huaweiEndpoint
}

// BuildUrl
// 构建华为的 url, baseUrl 必须以 “/”  开头
func (huaweiClient *HuaweiClient) BuildUrl(serviceType service_type.Type,
	baseUrl string, args interface{}, path ...string) *HuaweiClient {
	huaweiClient.url = huaweiClient.BuildEndpoint(serviceType, baseUrl, args, path...)
	// 打印请求的 url
	// klog.Info(fmt.Sprintf("huaweiClient.url: %v\n", huaweiClient.url))
	return huaweiClient
}

func (huaweiClient *HuaweiClient) BuildUrlWithProject(serviceType service_type.Type,
	baseUrl string, args interface{}, path ...string) *HuaweiClient {
	// 如果 baseUrl 包含{project_id} 则需要替换成真实的 projectId
	if strings.Contains(baseUrl, "{project_id}") {
		projectId, err := huaweiClient.GetProjectId()
		if err != nil {
			// todo 需要将错误抛出
			klog.Error("get project id err: ", err)
			huaweiClient.SetErr(errors.New(fmt.Sprint("get project id err: ", err)))
		}
		baseUrl = strings.Replace(baseUrl, "{project_id}", projectId, 1)
	}
	huaweiClient.url = huaweiClient.BuildEndpoint(serviceType, baseUrl, args, path...)
	// 打印请求的 url
	// klog.Info(fmt.Sprintf("huaweiClient.url: %v\n", huaweiClient.url))
	return huaweiClient
}

// Query
// 将 args 参数转换成 a=aValue&b=bValue url value 的格式
func Query(args interface{}) url.Values {
	v, err := query.Values(args, "huawei")
	if err != nil {
		klog.Error("encode query err: ", err)
	}
	return v
}

func BuildBaseUrl(baseUrl string, args interface{}, path ...string) string {
	url := baseUrl
	s := Query(args).Encode()
	for _, v := range path {
		start := strings.Index(url, "{")
		end := strings.Index(url, "}") + 1
		if start != -1 && end != -1 && start < end {
			url = string_utils.ReplaceAtPosition(url, v, start, end)
		}
	}
	if s == "" {
		return url
	}
	return fmt.Sprintf("%s?%s", url, s)
}

// MapToJson
// 将 url.values 转成 string
func MapToJson(param url.Values) string {
	if len(param) > 0 {
		dataType, _ := json.Marshal(param)
		dataString := string(dataType)
		return dataString
	}
	return ""
}
