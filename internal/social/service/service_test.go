package service

import (
	"douyin/api/pb/social"
	sdb "douyin/internal/social/db"
	"douyin/pkg/db"
	"douyin/pkg/logger"
	"fmt"
	"testing"
)

func TestNewService(t *testing.T) {
	// 初始化日志
	logger.Init(nil)
	sql := &db.Mysql{Username: "root", Password: "my password", Host: "my host", Port: 3306, Dbname: "social"}
	rds := &db.Redis{
		Host:     "127.0.0.1",
		Password: "",
		Port:     6379,
	}
	s := NewService(sdb.NewSocialDB(sql, rds))
	if s == nil {
		fmt.Printf("初始化错误")
		return
	}
	r, err :=
		//relation/action/ - 关系操作
		// s.RelationAction(&social.RpcRelationActionRequest{UserId: 6, ToUserId: 1, ActionType: 1}) //------测试通过
		//relatioin/follow/list/ - 关注列表
		// s.RelationFollowerList(&social.DouyinRelationFollowerListRequest{UserId: 6}) //------测试通过
		// relatioin/follow/list/ - 关注列表
		// s.RelationFollowList(&social.DouyinRelationFollowListRequest{UserId: 1}) //------测试通过
		//relation/friend/list/ - 好友列表
		// s.RelationFriendList(&social.DouyinRelationFriendListRequest{UserId: 1}) //------测试通过
		//message/send/ - 消息操作
		// s.MessageSend(&social.RpcMessageSendRequest{UserId: 1, ToUserId: 6, ActionType: 1, Content: "hello"}) //------测试通过
		// message/chat/ - 聊天记录
		s.MessageChat(&social.RpcMessageChatRequest{UserId: 1, ToUserId: 6}) //------测试通过
	if err != nil {
		fmt.Printf("测试失败,err:%v", err)
		return
	}
	fmt.Printf("测试通过  r:%v", r)

}
