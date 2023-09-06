package main

import (
	"douyin/internal/social/config"
	"douyin/internal/social/db"
	"douyin/internal/social/handler"
	"douyin/internal/social/service"
	etcdinit "douyin/pkg/etcd_init"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"

	"github.com/spf13/pflag"
)

var configPath = pflag.StringP("config", "c", "configs/social.yaml", "配置文件路径")

func main() {
	// 解析命令行参数
	pflag.Parse()

	// 读取配置文件
	config.Init(*configPath)

	// 初始化日志工具
	logger.Init(config.GlobalConfig.Log)

	// 打印配置文件测试
	logger.Debugf("服务=%#v\n", config.GlobalConfig.RPC.Social)

	// 初始化 grpc 服务
	s := grpc_server.NewGrpcServer(config.GlobalConfig.RPC.Social)

	//初始化所有依赖:
	// 初始化handler
	h := InitHandlers()
	//关闭数据库,rdies
	defer h.SocialHandlers.SocialService.Close()

	// 注册服务
	handler.Register(s, h.SocialHandlers)

	//关闭数据库,rdies
	defer h.SocialHandlers.SocialService.Close()

	etcdRegister := etcdinit.InitETCD("/etcd/social", config.GlobalConfig.RPC.Social, 5)
	//初始化etcd
	defer etcdRegister.Close()

	// 启动服务
	s.Start()
	defer s.Stop()
}

type Handlers struct {
	SocialHandlers *handler.SocialHandler
}

// 初始化所有依赖
func InitHandlers() *Handlers {
	// TODO: 初始化 gorm 和 redis

	// 初始化db, TODO: 依赖注入
	socialDB := db.NewSocialDB(&config.GlobalConfig.Mysql, &config.GlobalConfig.Redis)
	//初始化service
	socialService := service.NewService(socialDB)
	//初始化handler
	socialHandler := handler.NewSocialHandler(&socialService)

	return &Handlers{
		socialHandler,
	}
}
