package http

import (
	"context"
)

type BaseRequest struct {
	context.Context
	method string
	url    string
	body   interface{}
	header map[string][]string
}

func NewBaseRequest() *BaseRequest {
	return &BaseRequest{}
}

func (req *BaseRequest) Method(method string) *BaseRequest {
	req.method = method
	return req
}

func (req *BaseRequest) Url(url string) *BaseRequest {
	req.url = url
	return req
}

func (req *BaseRequest) Body(body interface{}) *BaseRequest {
	req.body = body
	return req
}

func (req *BaseRequest) Header(header map[string][]string) *BaseRequest {
	req.header = header
	return req
}

func (req *BaseRequest) SetMethod(method string) {
	req.method = method
}

func (req *BaseRequest) GetMethod() string {
	return req.method
}

func (req *BaseRequest) SetBody(body interface{}) {
	req.body = body
}

func (req *BaseRequest) GetBody() interface{} {
	return req.body
}

func (req *BaseRequest) SetUrl(url string) {
	req.url = url
}

func (req *BaseRequest) GetUrl() string {
	return req.url
}

func (req *BaseRequest) SetHeader(header map[string][]string) {
	req.header = header
}

func (req *BaseRequest) GetHeader() map[string][]string {
	return req.header
}
