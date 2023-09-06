package service

import (
	"context"
	"douyin/api/pb/social"
	"douyin/internal/gateway/db"
	"douyin/pkg/grpc_client"
	"douyin/pkg/logger"

	"github.com/gomodule/redigo/redis"
)

// ISocialClient Grpc 调用用户服务
type ISocialClient interface {
	//继承db的接口
	db.IGetWay
	//relation/action/ - 关系操作
	RelationAction(req *social.RpcRelationActionRequest) (resp *social.DouyinRelationActionResponse, err error)
	//relatioin/follow/list/ - 关注列表
	RelationFollowList(req *social.DouyinRelationFollowListRequest) (resp *social.DouyinRelationFollowListResponse, err error)
	//relation/follower/list/ - 粉丝列表
	RelationFollowerList(req *social.DouyinRelationFollowerListRequest) (resp *social.DouyinRelationFollowerListResponse, err error)
	//relation/friend/list/ - 好友列表
	RelationFriendList(req *social.DouyinRelationFriendListRequest) (resp *social.DouyinRelationFriendListResponse, err error)
	//message/chat/ - 聊天记录
	MessageChat(req *social.RpcMessageChatRequest) (resp *social.DouyinMessageChatResponse, err error)
	//message/send/ - 发送消息
	MessageSend(req *social.RpcMessageSendRequest) (resp *social.DouyinMessageSendResponse, err error)
}

// 定义结构体
type SocialClient struct {
	db.IGetWay
	SocialClient social.SocialClient
}

// 实例化
func NewSocialClient(pool *redis.Pool, conn *grpc_client.GrpcClient) ISocialClient {
	return &SocialClient{
		SocialClient: social.NewSocialClient(conn.ClientConn),
		IGetWay:      &db.GetWay{Pool: pool},
	}
}

// relation/action/ - 关系操作
func (s *SocialClient) RelationAction(req *social.RpcRelationActionRequest) (resp *social.DouyinRelationActionResponse, err error) {

	resp, err = s.SocialClient.DouyinRelationAction(context.Background(), req)
	if err != nil {
		logger.Errorf("关系操作relation/action: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// relatioin/follow/list/ - 关注列表
func (s *SocialClient) RelationFollowList(req *social.DouyinRelationFollowListRequest) (resp *social.DouyinRelationFollowListResponse, err error) {

	resp, err = s.SocialClient.DouyinRelationFollowList(context.Background(), req)
	if err != nil {
		logger.Errorf("关注列表relatioin/follow/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// relation/follower/list/ - 粉丝列表
func (s *SocialClient) RelationFollowerList(req *social.DouyinRelationFollowerListRequest) (resp *social.DouyinRelationFollowerListResponse, err error) {

	resp, err = s.SocialClient.DouyinRelationFollowerList(context.Background(), req)
	if err != nil {
		logger.Errorf("粉丝列表relation/follower/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// relation/friend/list/ - 好友列表
func (s *SocialClient) RelationFriendList(req *social.DouyinRelationFriendListRequest) (resp *social.DouyinRelationFriendListResponse, err error) {

	resp, err = s.SocialClient.DouyinRelationFriendList(context.Background(), req)
	if err != nil {
		logger.Errorf("好友列表relation/friend/list: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// message/chat/ - 聊天记录
func (s *SocialClient) MessageChat(req *social.RpcMessageChatRequest) (resp *social.DouyinMessageChatResponse, err error) {

	resp, err = s.SocialClient.DouyinMessageChat(context.Background(), req)
	if err != nil {
		logger.Errorf("聊天记录message/chat: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// message/send/ - 发送消息
func (s *SocialClient) MessageSend(req *social.RpcMessageSendRequest) (resp *social.DouyinMessageSendResponse, err error) {

	resp, err = s.SocialClient.DouyinMessageSend(context.Background(), req)
	if err != nil {
		logger.Errorf("发送消息message/send: %s\n", err)
		return nil, err
	}

	return resp, nil
}
