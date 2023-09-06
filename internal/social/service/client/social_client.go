package client

import (
	"context"
	"douyin/api/pb/social"
	discover "douyin/pkg/etcd_client"
	"douyin/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// 注册social rpc服务
// 定义接口
type ISocialClient interface {
	//批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map中
	FindFollows(req *social.FindFollowsRequest) (resp *social.FindFollowsResponse, err error)
}

// 定义结构体
type SocialClient struct {
	social.UnimplementedSocialServer

	SocialClient social.SocialClient
}

// 实现new方法
func NewSocialClient() ISocialClient {
	etcdResolverBuilder := discover.NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)
	// 初始化 grpc_client
	Conn, err := grpc.Dial("etcd:/etcd/social", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] 连接异常： %s\n", err)
		panic(err)
	}
	return &SocialClient{
		SocialClient: social.NewSocialClient(Conn),
	}
}

// 批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map中
func (s *SocialClient) FindFollows(req *social.FindFollowsRequest) (resp *social.FindFollowsResponse, err error) {
	resp, err = s.SocialClient.FindFollows(context.Background(), req)
	if err != nil {
		logger.Errorf("[error] 连接异常： %s\n", err)
		return nil, err
	}
	return resp, nil
}
