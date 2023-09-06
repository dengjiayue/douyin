package main

import (
	"douyin/internal/gateway/config"
	"douyin/internal/gateway/handler"
	middleware "douyin/internal/gateway/middleWare"
	"douyin/internal/gateway/service"
	"douyin/pkg/db"
	discover "douyin/pkg/etcd_client"
	"douyin/pkg/grpc_client"
	"douyin/pkg/logger"
	"douyin/pkg/web"
	"log"

	"google.golang.org/grpc"

	"github.com/spf13/pflag"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

var configPath = pflag.StringP("config", "c", "configs/gateway.yaml", "配置文件路径")

func main() {
	// 解析命令行参数
	pflag.Parse()

	// 读取配置文件
	config.Init(*configPath)
	log.Printf("配置文件成功")

	// 初始化日志工具
	logger.Init(config.GlobalConfig.Log)

	// 打印配置文件测试
	logger.Debugf("服务=%#v\n", "网关")
	logger.Debugf(config.GlobalConfig.RPC.User)

	// 初始化 Web 服务
	s := web.NewWeb(config.GlobalConfig.Web.Addr)
	//初始化web服务成功
	logger.Debugf("初始化web服务成功")

	// 初始化handler
	d := InitHandlers()
	//初始化hexitandler成功
	logger.Debugf("初始化handler成功")

	//注册中间件
	//输出响应数据
	s.Use(middleware.LogResponseDataMiddleware)
	//鉴权
	s.Use(middleware.AuthMiddlewareQueryOrForm(d.socialHandler.SocialClient))

	// 注册服务
	handler.UserGin(s, d.UserHandler)
	handler.VideoListGin(s, d.videoListHandler)
	handler.VideoGin(s, d.videoHandler)
	handler.SocialGin(s, d.socialHandler)

	//注册服务成功
	logger.Debugf("注册服务成功")

	//结束关闭redis连接池
	defer d.UserHandler.UserClient.Close()

	// 启动服务
	s.Start()
	defer s.Stop()

}

type Handlers struct {
	UserHandler      *handler.UserHandler
	videoListHandler *handler.VideoListHandler
	videoHandler     *handler.VideoHandler
	socialHandler    *handler.SocialHandler
}

// 初始化所有依赖
func InitHandlers() *Handlers {

	etcdResolverBuilder := discover.NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)
	// TODO: 初始化 gorm 和 redis

	// 初始化 grpc_client
	//user连接
	UConn, err := grpc.Dial("etcd:/etcd/user", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] user 连接异常： %s\n", err)
		panic(err)
	}
	//video连接
	VLConn, err := grpc.Dial("etcd:/etcd/video_list", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] video 连接异常： %s\n", err)
		panic(err)
	}
	//social连接
	SConn, err := grpc.Dial("etcd:/etcd/social", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] video 连接异常： %s\n", err)
		panic(err)
	}
	//video连接
	VConn, err := grpc.Dial("etcd:/etcd/video/", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] video 连接异常： %s\n", err)
		panic(err)
	}

	userConn := &grpc_client.GrpcClient{ClientConn: UConn}

	videoListConn := &grpc_client.GrpcClient{ClientConn: VLConn}

	socialConn := &grpc_client.GrpcClient{ClientConn: SConn}

	videoConn := &grpc_client.GrpcClient{ClientConn: VConn}

	//构建redis连接池
	pool := db.NewRedisPool(&db.Redis{Host: "127.0.0.1", Port: 6379, Password: "douyinapp"})

	// 初始化service
	userClient := service.NewUserClient(pool, userConn)

	videoListClient := service.NewVideoListClient(pool, videoListConn)

	socialClient := service.NewSocialClient(pool, socialConn)

	videoClient := service.NewVideoClient(pool, videoConn)

	// 初始化handler
	userHandler := handler.NewUserHandler(&userClient)
	videoListHandler := handler.NewVideoListHandler(&videoListClient)
	socialHandler := handler.NewSocialHandler(&socialClient)
	videoHandler := handler.NewVideoHandler(&videoClient)

	return &Handlers{
		UserHandler:      userHandler,
		videoListHandler: videoListHandler,
		socialHandler:    socialHandler,
		videoHandler:     videoHandler,
	}
}
