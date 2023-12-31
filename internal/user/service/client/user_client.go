package client

import (
	"context"
	"douyin/api/pb/user"
	discover "douyin/pkg/etcd_client"
	"douyin/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// 注册user rpc服务
// 定义接口
type IUserClient interface {
	//批量获取用户信息
	UsersInfo(req *user.DouyinUsersRequest) (resp *user.DouyinUsersResponse, err error)
	//获取用户信息
	UserInfo(req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error)
	//更新用户数据
	UserChange(req *user.DouyinUserChangeRequest) (resp *user.DouyinUserChangeResponse, err error)
}

// 定义结构体
type UserClient struct {
	user.UnimplementedUserServer

	UserClient user.UserClient
}

// 实现new方法
func NewUserClient() IUserClient {
	etcdResolverBuilder := discover.NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)
	// 初始化 grpc_client
	Conn, err := grpc.Dial("etcd:/etcd/user", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		logger.Errorf("[error] 连接异常： %s\n", err)
		panic(err)
	}
	return &UserClient{
		UserClient: user.NewUserClient(Conn),
	}
}

// 批量获取用户信息
func (s *UserClient) UsersInfo(req *user.DouyinUsersRequest) (resp *user.DouyinUsersResponse, err error) {
	resp, err = s.UserClient.DouyinUsersInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 获取用户信息
func (s *UserClient) UserInfo(req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {
	resp, err = s.UserClient.DouyinUserInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 更新用户数据
func (s *UserClient) UserChange(req *user.DouyinUserChangeRequest) (resp *user.DouyinUserChangeResponse, err error) {
	resp, err = s.UserClient.DouyinUserChange(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
