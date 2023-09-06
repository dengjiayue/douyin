package main

import (
	"douyin/internal/user/config"
	"douyin/internal/user/db"
	"douyin/internal/user/handler"
	"douyin/internal/user/service"
	etcdinit "douyin/pkg/etcd_init"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"

	"github.com/spf13/pflag"
)

var configPath = pflag.StringP("config", "c", "configs/user.yaml", "配置文件路径")

func main() {
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
	s := grpc_server.NewGrpcServer(config.GlobalConfig.RPC.User)

	// 初始化handler
	h := InitHandlers()
	//关闭数据库,rdies
	defer h.UserHandler.UserService.Close()

	// 注册服务
	handler.Register(s, h.UserHandler)

	//初始化etcd
	defer etcdinit.InitETCD("/etcd/user", config.GlobalConfig.RPC.User, 5).Close()

	// 启动服务
	s.Start()
	defer s.Stop()

}

type Handlers struct {
	UserHandler *handler.UserHandler
}

// 初始化所有依赖
func InitHandlers() *Handlers {
	// TODO: 初始化 gorm 和 redis

	// 初始化db, TODO: 依赖注入
	userDB := db.NewUserDB(&config.GlobalConfig.Mysql, &config.GlobalConfig.Redis)

	// 初始化service
	userService := service.NewUserService(userDB)

	// 初始化handler
	userHandler := handler.NewUserHandler(&userService)

	return &Handlers{
		UserHandler: userHandler,
	}
}
