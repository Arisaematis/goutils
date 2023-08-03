package cloud_provider

import (
	"fmt"
	"testing"
)

func TestCloudProviderId(t *testing.T) {
	var f Id = Huawei
	fmt.Println(f)
	switch f {
	case Huawei:
		fmt.Println("huawei cloudProviderId")
	case Hcso:
		fmt.Println("hcso cloudProviderId")
	case Hcs:
		fmt.Println("hcs cloudProviderId")
	default:
		fmt.Println("没有这样的服务商Id ...")
	}
}
