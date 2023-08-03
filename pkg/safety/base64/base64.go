package base64

import "encoding/base64"

func Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func Decode(dst []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(dst[:]))
}
