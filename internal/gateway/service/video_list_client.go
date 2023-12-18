package service

import (
	"context"
	"douyin/api/pb/video_list"
	"douyin/internal/gateway/db"
	"douyin/pkg/grpc_client"
	"douyin/pkg/logger"

	"github.com/gomodule/redigo/redis"
)

// IVideoListClient Grpc 调用用户服务
type IVideoListClient interface {
	//继承db的接口
	db.IGetWay
	// 批量查询视频信息
	VideoList(req *video_list.RpcFeedRequest) (resp *video_list.DouyinFeedResponse, err error)
}

type VideoListClient struct {
	db.IGetWay
	VideoListClient video_list.VideoListClient
}

func NewVideoListClient(pool *redis.Pool, conn *grpc_client.GrpcClient) IVideoListClient {
	return &VideoListClient{
		VideoListClient: video_list.NewVideoListClient(conn.ClientConn),
		IGetWay:         &db.GetWay{Pool: pool},
	}
}

// 批量查询视频信息
func (s *VideoListClient) VideoList(req *video_list.RpcFeedRequest) (resp *video_list.DouyinFeedResponse, err error) {

	resp, err = s.VideoListClient.DouyinFeed(context.Background(), req)
	if err != nil {
		logger.Errorf("批量查询视频信息: %s\n", err)
		return nil, err
	}

	return resp, nil
}
