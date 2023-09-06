package handler

import (
	"context"
	"douyin/api/pb/video_list"
	"douyin/internal/video_list/service"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"
)

// 注册服务
func Register(s *grpc_server.GrpcServer, vh *VideoListHandler) {
	logger.Debugf("注册 VideoServer")
	video_list.RegisterVideoListServer(s, vh)
	logger.Debugf("注册完成")
}

// 定义结构体
type VideoListHandler struct {
	video_list.UnimplementedVideoListServer
	//继承接口
	VideoListService service.IVideoService
}

// 实例化结构体
func NewVideoListHandler(VideoListService *service.IVideoService) *VideoListHandler {
	return &VideoListHandler{
		VideoListService: *VideoListService,
	}
}

// feed流查询视频信息
func (s *VideoListHandler) DouyinFeed(ctx context.Context, req *video_list.RpcFeedRequest) (*video_list.DouyinFeedResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoListService.Feed(req)
}
