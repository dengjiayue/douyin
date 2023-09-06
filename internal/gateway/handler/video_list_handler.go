package handler

import (
	"douyin/api/pb/video_list"
	"douyin/internal/gateway/service"
	"douyin/pkg/logger"
	"douyin/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin
func VideoListGin(s *web.Web, vh *VideoListHandler) {

	//获取视频列表
	s.GET("/douyin/feed/", vh.DouyinVideoList)
	//空url
	s.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status_code": "ok",
		})
	})
}

// handler
type VideoListHandler struct {
	video_list.UnimplementedVideoListServer
	//
	VideoListClient service.IVideoListClient
}

// NewVideoListHandler 实例化
func NewVideoListHandler(VideoListClient *service.IVideoListClient) *VideoListHandler {
	return &VideoListHandler{
		VideoListClient: *VideoListClient,
	}
}

// feed流查询视频信息
func (s *VideoListHandler) DouyinVideoList(c *gin.Context) {
	var req video_list.RpcFeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400,
			ErrorRespInt(400, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		req.UserId = 0
	} else {
		req.UserId = id
	}

	resp, err := s.VideoListClient.VideoList(&req)
	if err != nil {
		c.JSON(500,
			ErrorRespInt(500, "调用视频列表服务异常"))
		logger.Errorf("调用视频列表服务异常: %s\n", err)
		return
	}

	c.JSON(http.StatusOK,
		resp)
}
