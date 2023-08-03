package sdk_sign

import "net/http"

type ISign interface {
	Sign(r *http.Request) error
	GetSecretKey() string
	GetSecret() string
}

type DefaultSign struct {
}

func (d *DefaultSign) Sign(r *http.Request) error {
	return nil
}

func (d *DefaultSign) GetSecretKey() string {
	return ""
}
func (d *DefaultSign) GetSecret() string {
	return ""
}

type AliSigner struct {
}

func (a *AliSigner) Sign(r *http.Request) error {
	return nil
}
