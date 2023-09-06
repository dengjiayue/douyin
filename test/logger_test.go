package main

import (
	"douyin/pkg/logger"
	"testing"
	"time"
)

func TestConfigInit(t *testing.T) {
	logger.Init(nil)
	logger.Debugf("服务=%#v\n", "网关")
	time.Sleep(time.Second * 1)
}
