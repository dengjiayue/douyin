package main

import (
	"douyin/internal/video/config"
	"douyin/internal/video/db"
	"douyin/internal/video/handler"
	"douyin/internal/video/service"
	etcdinit "douyin/pkg/etcd_init"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"

	"github.com/spf13/pflag"
)

var configPath = pflag.StringP("config", "c", "configs/video.yaml", "配置文件路径")

func main() {
	// 解析命令行参数
	pflag.Parse()

	// 读取配置文件
	config.Init(*configPath)

	// 初始化日志工具
	logger.Init(config.GlobalConfig.Log)

	// 打印配置文件测试
	logger.Debugf("服务=%#v\n", config.GlobalConfig.RPC.Video)

	// 初始化 grpc 服务
	s := grpc_server.NewGrpcServer(config.GlobalConfig.RPC.Video)

	//初始化所有依赖
	h := InitHandlers()
	defer h.VideoHandlers.VideoService.Close()

	// 注册服务
	handler.Register(s, h.VideoHandlers)

	//初始化etcd
	etcdRegister := etcdinit.InitETCD("/etcd/video/", config.GlobalConfig.RPC.Video, 5)
	defer etcdRegister.Close()

	// 启动服务
	s.Start()
	defer s.Stop()
}

type Handlers struct {
	VideoHandlers *handler.VideoHandler
}

// 初始化所有依赖
func InitHandlers() *Handlers {
	// TODO: 初始化 gorm 和 redis

	// 初始化db, TODO: 依赖注入
	videoDB := db.NewVideoDB(&config.GlobalConfig.Mysql, &config.GlobalConfig.Redis)
	//初始化service
	videoService := service.NewService(videoDB)
	//初始化handler
	videoHandler := handler.NewVideoHandler(&videoService)

	return &Handlers{
		videoHandler,
	}

}
