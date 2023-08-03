package cloud_provider

import (
	"errors"
	factory "goutils/pkg/cloud/cloud_request"
	huaweiCloud "goutils/pkg/cloud/cloud_request/huawei"
)

/*
 * @author: yeshibo
 * @date: Tuesday, 2022/09/06, 4:47:10 pm
 */

type Id int

const (
	Huawei Id = iota
	Hcso
	Hcs
)

func (c Id) String() string {
	return [...]string{"huawei", "hcso", "hcs"}[c]
}

// GetRequestFactory
// 当新增云服务商时,需要新增云服务商的工厂
func (c Id) GetRequestFactory() factory.IRequestProducerFactory {
	factory := []factory.IRequestProducerFactory{
		huaweiCloud.GetHuaweiProducerFactoryInstance(),
		// todo 需要新增
		nil,
		nil,
	}
	if int(c) > len(factory) {
		return nil
	}
	return factory[c]
}

func GetCloudProviderId(cloudProviderId string) (Id, error) {
	m := map[string]Id{
		"huaweiyun": Huawei,
		"hcso":      Hcso,
		"hcs":       Hcs,
	}
	i, ok := m[cloudProviderId]
	if ok {
		return i, nil
	}
	return 0, errors.New("cloudProviderId not find")
}
