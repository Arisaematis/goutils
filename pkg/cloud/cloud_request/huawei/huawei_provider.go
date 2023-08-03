package huawei

import (
	factory "goutils/pkg/cloud/cloud_request"
	accessKey "goutils/pkg/cloud/entity/accesskey"
	"sync"
)

/*
 * @author: yeshibo
 * @date: Tuesday, 2022/09/06, 10:45:35 am
 */

var once sync.Once

type HuaweiRequestProducerFactory struct {
}

func GetHuaweiProducerFactoryInstance() *HuaweiRequestProducerFactory {
	var huaweiFactoryInstance *HuaweiRequestProducerFactory
	once.Do(func() {
		huaweiFactoryInstance = &HuaweiRequestProducerFactory{}
	})
	return huaweiFactoryInstance
}

func (huawei *HuaweiRequestProducerFactory) GetRegionRequestProducer(ack accessKey.AccessKey) factory.BaseRegionRequestProducer {
	return GetRegionInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetEcsRequestProducer(ack accessKey.AccessKey) factory.BaseEcsRequestProducer {
	return GetEcsInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetCESRequestProducer(ack accessKey.AccessKey) factory.BaseCesRequestProducer {
	return GetCesInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetEvsRequestProducer(ack accessKey.AccessKey) factory.BaseEvsRequestProducer {
	return GetEvsInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetObsRequestProducer(ack accessKey.AccessKey) factory.BaseObsRequestProducer {
	return GetObsInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetNasRequestProducer(ack accessKey.AccessKey) factory.BaseNasRequestProducer {
	return GetNasInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetEipRequestProducer(ack accessKey.AccessKey) factory.BaseEipRequestProducer {
	return GetEipInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetVpcRequestProducer(ack accessKey.AccessKey) factory.BaseVpcRequestProducer {
	return GetVpcInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetImsRequestProducer(ack accessKey.AccessKey) factory.BaseImsRequestProducer {
	return GetImsInstance(ack)
}

func (huawei *HuaweiRequestProducerFactory) GetRdsRequestProducer(ack accessKey.AccessKey) factory.BaseRdsRequestProducer {
	return GetRdsInstance(ack)
}
