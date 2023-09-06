package handler

import (
	"douyin/api/pb/social"
	"douyin/internal/gateway/service"
	"douyin/pkg/logger"
	"douyin/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin
func SocialGin(s *web.Web, sh *SocialHandler) {

	//relation/action/ - 关系操作-
	logger.Debugf("关系操作 :/douyin/relation/action")
	s.POST("/douyin/relation/action/", sh.DouyinRelationAction)
	//relatioin/follow/list/ - 关注列表-
	logger.Debugf("关注列表 :/douyin/relatioin/follow/list")
	s.GET("/douyin/relation/follow/list/", sh.DouyinRelationFollowList)
	//relation/follower/list/ - 粉丝列表
	logger.Debugf("粉丝列表 :/douyin/relation/follower/list")
	s.GET("/douyin/relation/follower/list/", sh.DouyinRelationFollowerList)
	//relation/friend/list/ - 好友列表
	logger.Debugf("好友列表 :/douyin/relation/friend/list")
	s.GET("/douyin/relation/friend/list/", sh.DouyinRelationFriendList)
	///message/chat/ - 聊天记录
	logger.Debugf("聊天记录 :/douyin/message/chat")
	s.GET("/douyin/message/chat/", sh.DouyinMessageChat)
	//message/action- 消息操作
	logger.Debugf("消息操作 :/douyin/message/action")
	s.POST("/douyin/message/action/", sh.DouyinMessageAction)

}

// handler
type SocialHandler struct {
	social.UnimplementedSocialServer
	//
	SocialClient service.ISocialClient
}

// NewSocialHandler 实例化
func NewSocialHandler(SocialClient *service.ISocialClient) *SocialHandler {
	return &SocialHandler{
		SocialClient: *SocialClient,
	}
}

// 关系操作
func (s *SocialHandler) DouyinRelationAction(c *gin.Context) {
	var req social.RpcRelationActionRequest
	// 读取Query参数
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
	resp, err := s.SocialClient.RelationAction(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用关系操作服务异常"))
		logger.Errorf("调用关系操作服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 关注列表
func (s *SocialHandler) DouyinRelationFollowList(c *gin.Context) {
	var req social.DouyinRelationFollowListRequest
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
	resp, err := s.SocialClient.RelationFollowList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespString("500", "调用关注列表服务异常"))
		logger.Errorf("调用关注列表服务异常: %s\n", err)
		return
	}
	// resp.StatusCode = "0"
	c.JSON(http.StatusOK, resp)
}

// 粉丝列表
func (s *SocialHandler) DouyinRelationFollowerList(c *gin.Context) {
	var req social.DouyinRelationFollowerListRequest
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

	resp, err := s.SocialClient.RelationFollowerList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespString("500", "调用粉丝列表服务异常"))
		logger.Errorf("调用粉丝列表服务异常: %s\n", err)
		return
	}
	// resp.StatusCode = "0"
	c.JSON(http.StatusOK, resp)
}

// 好友列表
func (s *SocialHandler) DouyinRelationFriendList(c *gin.Context) {
	var req social.DouyinRelationFriendListRequest
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

	resp, err := s.SocialClient.RelationFriendList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespString("500", "调用好友列表服务异常"))
		logger.Errorf("调用好友列表服务异常: %s\n", err)
		return
	}
	// resp.StatusCode = "0"
	c.JSON(http.StatusOK, resp)
}

// 聊天记录
func (s *SocialHandler) DouyinMessageChat(c *gin.Context) {
	var req social.RpcMessageChatRequest
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
	req.UserId = id
	resp, err := s.SocialClient.MessageChat(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespString("500", "调用聊天记录服务异常"))
		logger.Errorf("调用聊天记录服务异常: %s\n", err)
		return
	}
	resp.StatusCode = "0"
	c.JSON(http.StatusOK, resp)
}

// 消息操作
func (s *SocialHandler) DouyinMessageAction(c *gin.Context) {
	var req social.RpcMessageSendRequest
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
	resp, err := s.SocialClient.MessageSend(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			ErrorRespInt(http.StatusInternalServerError, "调用消息操作服务异常"))
		logger.Errorf("调用消息操作服务异常: %s\n", err)
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 统一错误返回方法(code为string类型)
func ErrorRespString(code string, msg string) gin.H {
	return gin.H{
		"status_code": code,
		"status_msg":  msg,
	}
}

// 统一错误返回方法(code为int类型)
func ErrorRespInt(code int, msg string) gin.H {
	return gin.H{
		"status_code": code,
		"status_msg":  msg,
	}
}

// 将获取鉴权结果int64
func GetUid(uid any, b bool) (r int64) {
	r = uid.(int64)
	if r == 0 || !b {
		logger.Errorf("鉴权失败")
		return 0
	}
	return r
}
