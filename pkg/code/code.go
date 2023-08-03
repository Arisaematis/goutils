package code

import (
	"net/http"

	"github.com/lstack-org/go-web-framework/pkg/code"
)

var (
	ResourceIdNotFoundError = code.ServiceCode{
		HttpCode:     http.StatusOK,
		BusinessCode: ResourceIdNotFound,
		EnglishMsg:   "ResourceId not exist",
		ChineseMsg:   "指定的资源Id 未找到",
	}
)
