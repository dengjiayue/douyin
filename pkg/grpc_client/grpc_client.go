package grpc_client

import (
	"douyin/pkg/logger"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	*grpc.ClientConn
}

// NewGrpcClient 创建gRPC客户端
func NewGrpcClient(addr string) (*GrpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("连接异常： %s\n", err)
		return nil, err
	}

	return &GrpcClient{
		ClientConn: conn,
	}, nil

}

// Start 启动服务
func (s *GrpcClient) Start() {

}

// Stop 停止服务
func (s *GrpcClient) Stop() {
	s.ClientConn.Close()
}
