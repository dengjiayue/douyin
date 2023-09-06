package handler

import (
	"context"
	"douyin/api/pb/video"
	"douyin/internal/video/service"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"
)

// 注册video服务
func Register(s *grpc_server.GrpcServer, vh *VideoHandler) {
	logger.Debugf("注册 VideoServer")
	video.RegisterVideoServer(s, vh)
}

type VideoHandler struct {
	video.UnimplementedVideoServer

	VideoService service.IVideoService
}

func NewVideoHandler(VideoService *service.IVideoService) *VideoHandler {
	return &VideoHandler{
		VideoService: *VideoService,
	}
}

func (s *VideoHandler) DouyinPublishAction(req video.Video_DouyinPublishActionServer) error {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	if err := s.VideoService.PublishAction(req); err != nil {
		logger.Errorf("上传视频publish/action: %s\n", err)
		return err
	}
	// req.SendAndClose(&video.DouyinPublishActionResponse{
	// 	StatusCode: 0,
	// 	StatusMsg:  "上传视频成功",
	// })
	logger.Debugf("上传视频成功")
	return nil
}

func (s *VideoHandler) DouyinPublishList(ctx context.Context, req *video.DouyinPublishListRequest) (*video.DouyinPublishListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoService.PublishList(req)
}

func (s *VideoHandler) DouyinFavoriteAction(ctx context.Context, req *video.RpcFavoriteActionRequest) (*video.DouyinFavoriteActionResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoService.FavoriteAction(req)
}

func (s *VideoHandler) DouyinFavoriteList(ctx context.Context, req *video.DouyinFavoriteListRequest) (*video.DouyinFavoriteListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoService.FavoriteList(req)
}

func (s *VideoHandler) DouyinCommentAction(ctx context.Context, req *video.RpcCommentActionRequest) (*video.DouyinCommentActionResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoService.CommentAction(req)
}

func (s *VideoHandler) DouyinCommentList(ctx context.Context, req *video.RpcCommentListRequest) (*video.DouyinCommentListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.VideoService.CommentList(req)
}
