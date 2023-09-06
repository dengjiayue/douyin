package grpc_server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	*grpc.Server
	listen net.Listener
}

// NewGrpcServer 创建gRPC服务
func NewGrpcServer(addr string) *GrpcServer {
	// 1.监听
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("监听异常:%s\n", err)
	}
	fmt.Printf("监听端口：%s\n", addr)
	// 2.实例化gRPC
	s := grpc.NewServer()

	return &GrpcServer{
		Server: s,
		listen: listener,
	}
}

// Start 启动服务
func (s *GrpcServer) Start() {
	s.Server.Serve(s.listen)
}

// Stop 停止服务
func (s *GrpcServer) Stop() {
	s.Server.Stop()
}
