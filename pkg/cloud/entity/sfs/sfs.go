package sfs

import (
	"goutils/pkg/cloud/entity"
	"goutils/pkg/cloud/entity/request"
)

type Nas struct {
	request.BaseCloudResponse
	Size      float64
	QuotaSize float64
}

type GetNasRequest struct {
	*entity.BaseCloudRequest
}
