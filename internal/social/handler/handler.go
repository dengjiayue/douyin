package handler

import (
	"context"
	"douyin/api/pb/social"
	"douyin/internal/social/service"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"
)

// 注册服务
func Register(s *grpc_server.GrpcServer, uh *SocialHandler) {
	logger.Debugf("注册 UserServer")
	social.RegisterSocialServer(s, uh)
}

type SocialHandler struct {
	social.UnimplementedSocialServer

	SocialService service.ISocialService
}

func NewSocialHandler(SocialService *service.ISocialService) *SocialHandler {
	return &SocialHandler{
		SocialService: *SocialService,
	}
}

// 接口实现小技巧
var _ social.SocialServer = (*SocialHandler)(nil)

// relation/action/ - 关系操作
func (s *SocialHandler) DouyinRelationAction(context context.Context, req *social.RpcRelationActionRequest) (*social.DouyinRelationActionResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.RelationAction(req)
}

// relatioin/follow/list/ - 关注列表
func (s *SocialHandler) DouyinRelationFollowList(context context.Context, req *social.DouyinRelationFollowListRequest) (*social.DouyinRelationFollowListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)

	return s.SocialService.RelationFollowList(req)
}

// relation/follower/list/ - 粉丝列表
func (s *SocialHandler) DouyinRelationFollowerList(context context.Context, req *social.DouyinRelationFollowerListRequest) (*social.DouyinRelationFollowerListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.RelationFollowerList(req)
}

// relation/friend/list/ - 好友列表
func (s *SocialHandler) DouyinRelationFriendList(context context.Context, req *social.DouyinRelationFriendListRequest) (*social.DouyinRelationFriendListResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.RelationFriendList(req)
}

// message/chat/ - 聊天记录
func (s *SocialHandler) DouyinMessageChat(context context.Context, req *social.RpcMessageChatRequest) (*social.DouyinMessageChatResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.MessageChat(req)
}

// message/send/ - 消息操作
func (s *SocialHandler) DouyinMessageSend(context context.Context, req *social.RpcMessageSendRequest) (*social.DouyinMessageSendResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.MessageSend(req)
}

// 批量获取用户关注的用户
func (s *SocialHandler) FindFollows(context context.Context, req *social.FindFollowsRequest) (*social.FindFollowsResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.SocialService.FindFollowsService(req)
}
