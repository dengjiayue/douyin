package service

import (
	"douyin/api/pb/social"
	"douyin/api/pb/user"
	"douyin/internal/social/db"
	"douyin/internal/social/model"
	"douyin/internal/user/service/client"
	"douyin/pkg/logger"
	"time"
)

// 定义ISocialService接口
type ISocialService interface {
	//继承db.ICommentDB接口
	db.ISocialDB
	//relation/action/ - 关系操作
	RelationAction(req *social.RpcRelationActionRequest) (resp *social.DouyinRelationActionResponse, err error)
	//relatioin/follow/list/ - 关注列表
	RelationFollowList(req *social.DouyinRelationFollowListRequest) (resp *social.DouyinRelationFollowListResponse, err error)
	//relation/follower/list/ - 粉丝列表
	RelationFollowerList(req *social.DouyinRelationFollowerListRequest) (resp *social.DouyinRelationFollowerListResponse, err error)
	//relation/friend/list/ - 好友列表
	RelationFriendList(req *social.DouyinRelationFriendListRequest) (resp *social.DouyinRelationFriendListResponse, err error)
	///message/chat/ - 聊天记录
	MessageChat(req *social.RpcMessageChatRequest) (resp *social.DouyinMessageChatResponse, err error)
	//message/send/ - 消息操作
	MessageSend(req *social.RpcMessageSendRequest) (resp *social.DouyinMessageSendResponse, err error)
	FindFollowsService(req *social.FindFollowsRequest) (*social.FindFollowsResponse, error)
}

// 定义结构体
type SocialService struct {
	//继承db.SocialDB
	db.ISocialDB
	//继承user_client接口
	client.IUserClient
}

// 确保结构体实现接口的小技巧:
var _ ISocialService = (*SocialService)(nil)

// 实例化
func NewService(socialDB db.ISocialDB) ISocialService {

	return &SocialService{socialDB, client.NewUserClient()}
}

// relation/action/ - 关系操作
func (s *SocialService) RelationAction(req *social.RpcRelationActionRequest) (*social.DouyinRelationActionResponse, error) {
	//更新关注表数据
	n, err := s.UpdateFollow(&model.Follow{Follower: req.UserId, Followed: req.ToUserId, IsFollow: req.ActionType})
	if err != nil {
		return &social.DouyinRelationActionResponse{StatusCode: 500, StatusMsg: "服务器错误!!!关注失败"}, err
	}
	if n == 0 {
		return &social.DouyinRelationActionResponse{StatusCode: 500, StatusMsg: "无法重复关注"}, err
	}

	//调用user服务更新用户数据
	_, err = s.UserChange(&user.DouyinUserChangeRequest{UserId: req.UserId, ToUserId: req.ToUserId, ActionType: req.ActionType, Type: 2})
	if err != nil {
		logger.Debugf("调用user服务更新用户数据失败 err:%v", err)
		return &social.DouyinRelationActionResponse{StatusCode: 500, StatusMsg: "服务器错误!!!关注失败"}, err
	}
	return &social.DouyinRelationActionResponse{StatusCode: 0, StatusMsg: "ok"}, nil
}

// relatioin/follow/list/ - 关注列表
func (s *SocialService) RelationFollowList(req *social.DouyinRelationFollowListRequest) (resp *social.DouyinRelationFollowListResponse, err error) {
	//根据关注者查询被关注者s
	followeds, err := s.GetFolloweds(req.UserId)
	if err != nil {
		return &social.DouyinRelationFollowListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询关注列表失败"}, err
	}
	//调用user服务查询查询用户信息s
	userList, err := s.UsersInfo(&user.DouyinUsersRequest{UserIds: followeds})
	if err != nil {
		logger.Errorf("调用user服务查询用户信息失败 err:%v", err)
		return &social.DouyinRelationFollowListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询关注列表失败"}, err
	}
	for i := 0; i < len(userList.Users); i++ {
		userList.Users[i].IsFollow = true
	}
	//返回数据
	resp = &social.DouyinRelationFollowListResponse{StatusCode: "0", StatusMsg: "ok", UserList: userList.Users}

	return
}

// relation/follower/list/ - 粉丝列表
func (s *SocialService) RelationFollowerList(req *social.DouyinRelationFollowerListRequest) (resp *social.DouyinRelationFollowerListResponse, err error) {
	//根据被关注者查询关注者s
	followers, err := s.GetFollowers(req.UserId)
	if err != nil {
		return &social.DouyinRelationFollowerListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询粉丝列表失败"}, err
	}
	//调用user服务查询查询用户信息s
	userList, err := s.UsersInfo(&user.DouyinUsersRequest{UserIds: followers})
	if err != nil {
		logger.Errorf("调用user服务查询用户信息失败 err:%v", err)
		return &social.DouyinRelationFollowerListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询粉丝列表失败"}, err
	}
	// 批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map
	follows, err := s.FindFollowsService(&social.FindFollowsRequest{UserId: req.UserId, Ids: followers})
	if err != nil {
		logger.Errorf("批量获取用户关注的用户失败 err:%v", err)
		return &social.DouyinRelationFollowerListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询粉丝列表失败"}, err
	}
	//遍历userList.Users,将follows中的数据添加到userList.Users中
	for i := 0; i < len(userList.Users); i++ {
		userList.Users[i].IsFollow = follows.FollowMap[userList.Users[i].Id]
	}
	//返回数据
	resp = &social.DouyinRelationFollowerListResponse{StatusCode: "0", StatusMsg: "ok", UserList: userList.Users}
	return
}

// relation/friend/list/ - 好友列表
func (s *SocialService) RelationFriendList(req *social.DouyinRelationFriendListRequest) (resp *social.DouyinRelationFriendListResponse, err error) {
	//根据关注者查询相互关注的用户(与关注者互相关注的用户)--friend
	friends, err := s.GetFriends(req.UserId)
	if err != nil {
		return &social.DouyinRelationFriendListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询好友列表失败"}, err
	}

	if len(friends) == 0 {
		return &social.DouyinRelationFriendListResponse{StatusCode: "0", StatusMsg: "没有好友", UserList: []*social.FriendUser{}}, nil
	}
	//调用user服务查询查询用户信息s
	users, err := s.UsersInfo(&user.DouyinUsersRequest{UserIds: friends})
	if err != nil {
		logger.Errorf("调用user服务查询用户信息失败 err:%v", err)
		return &social.DouyinRelationFriendListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询好友列表失败"}, err
	}

	//定义friend map
	friend_list := make(map[int64]*social.FriendUser)
	for _, user := range users.Users {
		friend_list[user.Id] = model.UserToFriend(user)
	}
	//查询roomids
	roomIds, err := s.GetRoomIds(req.UserId, friends)
	if err != nil {
		logger.Errorf("查询roomids失败 err:%v", err)
		return &social.DouyinRelationFriendListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询好友列表失败"}, err
	}
	if len(roomIds) != 0 {

		//查询与好友的最新一条聊天记录
		chatList, err := s.GetNewMessages(roomIds)
		if err != nil {
			logger.Errorf("查询与好友的最新一条聊天记录失败 err:%v", err)
			return &social.DouyinRelationFriendListResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询好友列表失败"}, err
		}
		//遍历chatList,将chatList中的数据添加到friend_list中(messagetype:0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息(req.UserId=message.FromUserId=>1;req.UserId=message.ToUserId=>0))
		for i := 0; i < len(chatList); i++ {
			//判断当前请求用户是发送者还是接收者
			if req.UserId == chatList[i].FromUserId {
				//当前请求用户是发送者
				friend_list[chatList[i].ToUserId].Message = chatList[i].Content
				friend_list[chatList[i].ToUserId].MsgType = 1
			} else {
				//当前请求用户是接收者
				friend_list[chatList[i].FromUserId].Message = chatList[i].Content
				friend_list[chatList[i].FromUserId].MsgType = 0
			}
		}
	}
	//遍历friend_list,将friend_list中的数据添加到friendList中
	friendList := make([]*social.FriendUser, len(friend_list))
	i := 0
	for _, friend := range friend_list {
		friendList[i] = friend
		i++
	}

	//返回数据
	resp = &social.DouyinRelationFriendListResponse{StatusCode: "0", StatusMsg: "ok", UserList: friendList}
	return
}

// message/chat/ - 聊天记录
func (s *SocialService) MessageChat(req *social.RpcMessageChatRequest) (resp *social.DouyinMessageChatResponse, err error) {
	//查询聊天记录
	messages, err := s.GetMessages(req.UserId, req.ToUserId, req.PreMsgTime)
	if err != nil {
		logger.Errorf("查询聊天记录失败 err:%v", err)
		return &social.DouyinMessageChatResponse{StatusCode: "500", StatusMsg: "服务器错误!!!查询聊天记录失败"}, err
	}
	//返回数据
	resp = &social.DouyinMessageChatResponse{StatusCode: "0", StatusMsg: "ok", MessageList: messages}

	return
}

// message/send/ - 消息操作
func (s *SocialService) MessageSend(req *social.RpcMessageSendRequest) (resp *social.DouyinMessageSendResponse, err error) {
	//查询/创建room
	room, err := s.GetRoom(req.UserId, req.ToUserId)
	if err != nil {
		logger.Errorf("查询/创建room失败 err:%v", err)
		return &social.DouyinMessageSendResponse{StatusCode: 500, StatusMsg: "服务器错误!!!发送消息失败"}, err
	}
	//插入message
	err = s.InsertMessage(&model.Message{Message: social.Message{FromUserId: req.UserId, ToUserId: req.ToUserId, Content: req.Content, CreateTime: time.Now().UnixMilli()}, RoomId: room.RoomId})
	if err != nil {
		logger.Errorf("插入message失败 err:%v", err)
		return &social.DouyinMessageSendResponse{StatusCode: 500, StatusMsg: "服务器错误!!!发送消息失败"}, err
	}
	//返回数据
	resp = &social.DouyinMessageSendResponse{StatusCode: 0, StatusMsg: "ok"}

	return
}

// 批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map
func (s *SocialService) FindFollowsService(req *social.FindFollowsRequest) (*social.FindFollowsResponse, error) {
	//根据ids查询用户关注的用户ids
	follows, err := s.GetFollows(req)
	if err != nil {
		logger.Errorf("查询用户关注的用户ids失败 err:%v", err)
		return nil, err
	}
	//返回数据
	return &social.FindFollowsResponse{FollowMap: follows}, nil

}
