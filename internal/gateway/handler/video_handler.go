package handler

import (
	"douyin/api/pb/video"
	"douyin/internal/gateway/service"
	"douyin/pkg/logger"
	"douyin/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin
func VideoGin(s *web.Web, vh *VideoHandler) {

	//上传视频publish/action
	logger.Debugf("上传视频 :/douyin/publish/action/")
	s.POST("/douyin/publish/action/", vh.DouyinPublishAction)
	//获取视频发布列表publish/list
	logger.Debugf("获取视频发布列表 :/douyin/publish/list")
	s.GET("/douyin/publish/list/", vh.DouyinPublishList)
	//点赞视频favorite/action
	logger.Debugf("点赞视频 :/douyin/favorite/action")
	s.POST("/douyin/favorite/action/", vh.DouyinFavoriteAction)
	//获取视频点赞列表favorite/list
	logger.Debugf("获取视频点赞列表 :/douyin/favorite/list")
	s.GET("/douyin/favorite/list/", vh.DouyinFavoriteList)
	//评论操作comment/action
	logger.Debugf("评论操作 :/douyin/comment/action")
	s.POST("/douyin/comment/action/", vh.DouyinCommentAction)
	//获取视频评论列表comment/list
	logger.Debugf("获取视频评论列表 :/douyin/comment/list")
	s.GET("/douyin/comment/list/", vh.DouyinCommentList)

}

// handler
type VideoHandler struct {
	video.UnimplementedVideoServer
	//
	VideoClient service.IVideoClient
}

// NewVideoHandler 实例化
func NewVideoHandler(VideoClient *service.IVideoClient) *VideoHandler {
	return &VideoHandler{
		VideoClient: *VideoClient,
	}
}

// 上传视频 publish/action
func (s *VideoHandler) DouyinPublishAction(c *gin.Context) {
	var req video.RpcPublishActionRequest

	// 校验参数
	req.Title = c.PostForm("title")
	logger.Debugf("title: %s\n", req.Title)

	//使用multipartForm获取参数--------------------------------------------弃用
	// form, err := c.MultipartForm()
	// if err != nil {
	// 	logger.Errorf("参数异常: %s\n", err)
	// 	return
	// }
	// data := form.File["data"][0]
	// formtoken := form.Value["token"][0]
	// logger.Debugf("formtoken: %s\n", formtoken)
	// logger.Debugf("data: %s\n", data.Filename)
	// title := form.Value["title"][0]
	// logger.Debugf("title: %s\n", title)
	// data, err = c.FormFile("data")
	// //校验参数
	// req.Title = title

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespInt(http.StatusUnauthorized, "鉴权失败"))
		return
	}
	req.UserId = id

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespInt(http.StatusBadRequest, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	resp, err := s.VideoClient.PublishAction(&req, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	logger.Debugf("ok")
	c.JSON(http.StatusOK, resp)
}

// 获取视频发布列表 publish/list
func (s *VideoHandler) DouyinPublishList(c *gin.Context) {
	var req video.DouyinPublishListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespInt(http.StatusBadRequest, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespInt(http.StatusUnauthorized, "鉴权失败"))
		return
	}
	req.Token = ""
	resp, err := s.VideoClient.PublishList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 点赞视频 favorite/action
func (s *VideoHandler) DouyinFavoriteAction(c *gin.Context) {
	var req video.RpcFavoriteActionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespInt(http.StatusBadRequest, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespInt(http.StatusUnauthorized, "鉴权失败"))
		return
	}
	req.UserId = id
	resp, err := s.VideoClient.FavoriteAction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 获取视频点赞列表 favorite/list-
func (s *VideoHandler) DouyinFavoriteList(c *gin.Context) {
	var req video.DouyinFavoriteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespString("400", "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespString("401", "鉴权失败"))
		return
	}
	req.Token = ""
	resp, err := s.VideoClient.FavoriteList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespString("500", "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	// resp.StatusCode = "0"
	c.JSON(http.StatusOK, resp)
}

// 评论操作 comment/action-
func (s *VideoHandler) DouyinCommentAction(c *gin.Context) {
	var req video.RpcCommentActionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespInt(http.StatusBadRequest, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}
	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespInt(http.StatusUnauthorized, "鉴权失败"))
		return
	}
	req.UserId = id
	resp, err := s.VideoClient.CommentAction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 获取视频评论列表 comment/list-
func (s *VideoHandler) DouyinCommentList(c *gin.Context) {
	var req video.RpcCommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			ErrorRespInt(http.StatusBadRequest, "参数异常"))
		logger.Errorf("参数异常: %s\n", err)
		return
	}

	//鉴权结果
	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("鉴权失败")
		c.JSON(http.StatusUnauthorized,
			ErrorRespInt(http.StatusUnauthorized, "鉴权失败"))
		return
	}
	req.UserId = id
	resp, err := s.VideoClient.CommentList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用视频服务异常"))
		logger.Errorf("调用视频服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}
