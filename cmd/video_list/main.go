package main

import (
	"douyin/internal/video_list/config"
	"douyin/internal/video_list/db"
	"douyin/internal/video_list/handler"
	"douyin/internal/video_list/service"
	etcdinit "douyin/pkg/etcd_init"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"

	"github.com/spf13/pflag"
)

// 暂时用不了,有bug,明天修复
func main() {
	var configPath = pflag.StringP("config", "c", "configs/video_list.yaml", "配置文件路径")

	// 解析命令行参数
	pflag.Parse()

	// 读取配置文件
	config.Init(*configPath)

	// 初始化日志工具
	logger.Init(config.GlobalConfig.Log)

	// 打印配置文件测试
	logger.Debugf("服务=%#v\n", "用户")
	logger.Debugf(config.GlobalConfig.RPC.User)

	// 初始化 grpc 服务
	s := grpc_server.NewGrpcServer(config.GlobalConfig.RPC.VideoList)

	//初始化etcd
	etcdRegister := etcdinit.InitETCD("/etcd/video_list", config.GlobalConfig.RPC.VideoList, 5)

	defer etcdRegister.Close()
	// 初始化handler
	d := InitHandlers()

	//关闭数据库,rdies
	defer d.VideoListHandler.VideoListService.Close()
	// 注册服务
	handler.Register(s, d.VideoListHandler)
	// 启动服务
	s.Start()
	defer s.Stop()
}

// 定义结构体
type Handlers struct {
	VideoListHandler *handler.VideoListHandler
}

// 所有依赖的服务都在这里初始化
func InitHandlers() *Handlers {
	//初始化db
	videodb := db.NewVideoListDB(&config.GlobalConfig.Mysql, &config.GlobalConfig.Redis)
	//初始化service
	videoservice := service.NewVideoListService(videodb)
	//初始化handler
	videohandler := handler.NewVideoListHandler(&videoservice)
	return &Handlers{
		VideoListHandler: videohandler,
	}
}
