package cron

import "github.com/robfig/cron/v3"

var Cron *cron.Cron

func NewWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// StartCron
// 开启定时任务
func init() {
	Cron = NewWithSeconds()
	Cron.Start()
}
