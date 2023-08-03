package download

import (
	"log"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	url := "https://lstack-image.oss-cn-hangzhou.aliyuncs.com/images/aaa.yaml?versionId=CAEQGRiBgIC_jZqi2RciIDdiMWE4MGM2NzZhMTQ4YTM4MmUyZTJhMWMwYWQwMDU2"
	fileName := "aaa.yaml"
	start := time.Now()
	d := NewDownloader(url, fileName, "/Users/yeshibo/Desktop/program/Go/goutils/pkg/excel", 3)
	if err := d.Download(); err != nil {
		log.Fatal(err)
	}
	log.Printf("下载%v 耗时：%v s", fileName, time.Since(start).Seconds())
}
