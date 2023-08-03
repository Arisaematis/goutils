package query

import (
	"fmt"
	"testing"
)

type Options struct {
	Query   string `huawei:"q,omitempty"`
	ShowAll bool   `huawei:"all,omitempty"`
	Page    int    `huawei:"page,omitempty"`
}

func TestQuery(t *testing.T) {
	// opt := Options{"foo", true, 2}
	opt := Options{
		Query:   "foo",
		ShowAll: false,
		Page:    0,
	}
	query(opt)
}

func query(args interface{}) {
	v, _ := Values(args, "huawei")
	fmt.Println(v.Encode())
	fmt.Println(v.Encode() == "")
}
