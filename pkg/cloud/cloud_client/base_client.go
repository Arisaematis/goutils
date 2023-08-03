package cloud_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	sdkSign "goutils/pkg/cloud/sdk_sign"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

type IBaseClient interface {
	SendRequest() ([]byte, error)
}

type ResponseErrorHandler func([]byte) error

type BaseClient struct {
	*http.Client
	context.Context
	Err    error
	Method string
	url    string
	body   []byte
	header map[string][]string
	sdkSign.ISign
	ResponseErrorHandler
}

func (base *BaseClient) GetHeader() map[string][]string {
	return base.header
}

// NewHttpClientIgnoreTLS 创建一个http client 忽略证书校验
func NewHttpClientIgnoreTLS() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		DisableKeepAlives:  false,
	}
	return &http.Client{Transport: tr}
}

// NewDefaultBaseClientWithContext 创建一个默认不带任何签名的client
func NewDefaultBaseClientWithContext(ctx context.Context) *BaseClient {
	return &BaseClient{
		Client:  NewHttpClientIgnoreTLS(),
		Context: ctx,
		ISign:   &sdkSign.DefaultSign{},
	}
}

// NewBaseClientWithContext 创建一个的 AK/SK 签名的客户端
func NewBaseClientWithContext(ctx context.Context, sign sdkSign.ISign) *BaseClient {
	return &BaseClient{
		Client:  NewHttpClientIgnoreTLS(),
		Context: ctx,
		ISign:   sign,
	}
}

// SendRequest 发送Http请求
func (base *BaseClient) SendRequest() ([]byte, error) {
	if base.Err != nil {
		return nil, base.Err
	}
	start := time.Now().UTC().UnixMilli()
	r, err := http.NewRequestWithContext(base.Context,
		base.Method,
		base.url,
		ioutil.NopCloser(bytes.NewBuffer(base.body)))
	if err != nil {
		return nil, err
	}
	err = base.Sign(r)
	if err != nil {
		return nil, err
	}
	for k, v := range base.header {
		r.Header.Set(k, v[0])
	}
	resp, err := base.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if base.ResponseErrorHandler != nil {
		if err = base.ResponseErrorHandler(body); err != nil {
			return nil, err
		}
	}
	end := time.Now().UTC().UnixMilli()
	klog.V(6).Info(fmt.Sprintf("The time required header[%v]", r.Header))
	klog.V(6).Info(fmt.Sprintf("The time required url[%s] to complete the request : %v ms", base.url, end-start))
	// 打印请求返回体 上线后需要删除 该方法用于调试
	klog.V(6).Info(fmt.Sprintf("cloud response: %s", string(body)))
	return body, nil
}

func (base *BaseClient) Get() *BaseClient {
	base.Method = http.MethodGet
	return base
}

func (base *BaseClient) Post() *BaseClient {
	base.Method = http.MethodPost
	return base
}

func (base *BaseClient) Delete() *BaseClient {
	base.Method = http.MethodDelete
	return base
}

func (base *BaseClient) Put() *BaseClient {
	base.Method = http.MethodPut
	return base
}

func (base *BaseClient) Patch() *BaseClient {
	base.Method = http.MethodPatch
	return base
}

func (base *BaseClient) Body(body []byte) *BaseClient {
	base.body = body
	return base
}

func (base *BaseClient) SetBody(body interface{}) *BaseClient {
	b, _ := json.Marshal(body)
	base.body = b
	return base
}

func (base *BaseClient) Url(url string) *BaseClient {
	base.url = url
	return base
}

func (base *BaseClient) Header(header map[string][]string) *BaseClient {
	base.header = header
	return base
}

// 设置error
func (base *BaseClient) SetErr(err error) *BaseClient {
	base.Err = err
	return base
}

// 设置错误处理函数
func (base *BaseClient) SetResponseErrorHandler(handler ResponseErrorHandler) *BaseClient {
	base.ResponseErrorHandler = handler
	return base
}
