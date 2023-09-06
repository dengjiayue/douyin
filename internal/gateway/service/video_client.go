package service

import (
	"context"
	"douyin/api/pb/video"
	"douyin/internal/gateway/db"
	"douyin/pkg/grpc_client"
	"douyin/pkg/logger"
	"io"
	"mime/multipart"

	"github.com/gomodule/redigo/redis"
)

// IVideoClient Grpc 调用用户服务
type IVideoClient interface {
	//继承db的接口
	db.IGetWay
	// 上传视频publish/action
	PublishAction(req *video.RpcPublishActionRequest, fileHeader *multipart.FileHeader) (resp *video.DouyinPublishActionResponse, err error)
	// 获取视频发布列表publish/list
	PublishList(req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error)
	// 点赞视频favorite/action
	FavoriteAction(req *video.RpcFavoriteActionRequest) (resp *video.DouyinFavoriteActionResponse, err error)
	// 获取视频点赞列表favorite/list
	FavoriteList(req *video.DouyinFavoriteListRequest) (resp *video.DouyinFavoriteListResponse, err error)
	// 评论视频comment/action
	CommentAction(req *video.RpcCommentActionRequest) (resp *video.DouyinCommentActionResponse, err error)
	// 获取视频评论列表comment/list
	CommentList(req *video.RpcCommentListRequest) (resp *video.DouyinCommentListResponse, err error)
}

// 定义结构体
type VideoClient struct {
	db.IGetWay
	VideoClient video.VideoClient
}

// 实例化
func NewVideoClient(pool *redis.Pool, conn *grpc_client.GrpcClient) IVideoClient {
	return &VideoClient{
		VideoClient: video.NewVideoClient(conn.ClientConn),
		IGetWay:     &db.GetWay{Pool: pool},
	}
}

// 上传视频publish/action(流式上传视频数据:发送到rpc服务端)
func (s *VideoClient) PublishAction(req *video.RpcPublishActionRequest, fileHeader *multipart.FileHeader) (resp *video.DouyinPublishActionResponse, err error) {
	file, err := fileHeader.Open()
	stream, err := s.VideoClient.DouyinPublishAction(context.Background())
	if err != nil {
		logger.Errorf("上传视频publish/action: %s\n", err)
		return &video.DouyinPublishActionResponse{
			StatusCode: 500,
			StatusMsg:  "上传视频文件失败",
		}, err
	}
	if err != nil {
		logger.Errorf("上传视频publish/action: %s\n", err)
		return &video.DouyinPublishActionResponse{
			StatusCode: 500,
			StatusMsg:  "上传视频文件失败",
		}, err
	}
	defer file.Close()
	// 发送标题数据
	err = stream.Send(req)
	//缓冲区,每次读取1M
	buffer := make([]byte, 1024*1024)
	// 发送视频数据
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("上传视频publish/action: %s\n", err)
			return &video.DouyinPublishActionResponse{
				StatusCode: 500,
				StatusMsg:  "上传视频文件失败",
			}, err
		}
		//发送数据
		err = stream.Send(&video.RpcPublishActionRequest{
			Data: buffer[:bytesRead],
		})
		if err != nil {
			logger.Errorf("上传视频publish/action: %s\n", err)
			return &video.DouyinPublishActionResponse{
				StatusCode: 500,
				StatusMsg:  "上传视频文件失败",
			}, err
		}
	}

	// 接收rpc响应
	_, err = stream.CloseAndRecv()
	if err != nil {
		//如果流程正常结束,则返回nil
		if err == io.EOF {
			logger.Debugf("上传视频成功")
		} else {
			logger.Errorf("上传视频publish/action: %s\n", err)
			return &video.DouyinPublishActionResponse{
				StatusCode: 500,
				StatusMsg:  "上传视频文件失败",
			}, err
		}
	}
	return &video.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "上传视频成功",
	}, nil
}

// 获取视频发布列表publish/list
func (s *VideoClient) PublishList(req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error) {

	resp, err = s.VideoClient.DouyinPublishList(context.Background(), req)
	if err != nil {
		logger.Errorf("获取视频发布列表publish/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 点赞视频favorite/action
func (s *VideoClient) FavoriteAction(req *video.RpcFavoriteActionRequest) (resp *video.DouyinFavoriteActionResponse, err error) {

	//请求rpc参数
	logger.Debugf("点赞视频favorite/action: %v\n", req)
	resp, err = s.VideoClient.DouyinFavoriteAction(context.Background(), req)
	if err != nil {
		logger.Errorf("点赞视频favorite/action: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取视频点赞列表favorite/list
func (s *VideoClient) FavoriteList(req *video.DouyinFavoriteListRequest) (resp *video.DouyinFavoriteListResponse, err error) {

	resp, err = s.VideoClient.DouyinFavoriteList(context.Background(), &video.DouyinFavoriteListRequest{UserId: req.UserId})
	if err != nil {
		logger.Errorf("获取视频点赞列表favorite/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 评论视频comment/action
func (s *VideoClient) CommentAction(req *video.RpcCommentActionRequest) (resp *video.DouyinCommentActionResponse, err error) {

	resp, err = s.VideoClient.DouyinCommentAction(context.Background(), req)
	if err != nil {
		logger.Errorf("评论视频comment/action: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取视频评论列表comment/list
func (s *VideoClient) CommentList(req *video.RpcCommentListRequest) (resp *video.DouyinCommentListResponse, err error) {

	resp, err = s.VideoClient.DouyinCommentList(context.Background(), req)
	if err != nil {
		logger.Errorf("获取视频评论列表comment/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}
