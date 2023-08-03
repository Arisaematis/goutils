package accesskey

import "fmt"

type AccessKey struct {
	AckUid          string `json:"ackUid" bson:"id"`
	CloudProviderId string `json:"cloudProviderId" bson:"cloudProviderId"`
	AccessKeyId     string `json:"accesskeyId" bson:"accessKeyId"`
	Secret          string `json:"accesskeySecret" bson:"accessKeySecret"`
	AccessKeyType   int64  `json:"accessKeyType" bson:"accessKeyType"`
}

func (accessKey AccessKey) ToString() string {
	return fmt.Sprintf("AccessKey[ackUid = %s, cloudProviderId = %s, accesskeyId = %s]",
		accessKey.AckUid, accessKey.CloudProviderId, accessKey.AccessKeyId)
}
