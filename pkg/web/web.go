package web

import (
	"douyin/pkg/logger"

	"github.com/gin-gonic/gin"
)

// 初始化 gin Engine
type Web struct {
	*gin.Engine
	addr string
}

func NewWeb(addr string) *Web {
	return &Web{
		Engine: gin.Default(),
		addr:   addr,
	}
}

func (s *Web) Start() {
	logger.Debugf("web服务启动,监听地址:%s\n", s.addr)
	s.Run(s.addr)
}

func (s *Web) Stop() {
}
