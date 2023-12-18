package db

import (
	"douyin/api/pb/social"
	"douyin/internal/social/model"
	"douyin/pkg/db"
	"douyin/pkg/logger"
	"fmt"

	jsoner "github.com/json-iterator/go"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// 全局: key为roomid:userid,值为上次获取到最后一条消息的时间戳-----------------------------------------------------------------------------------------------弃用
// var getMessageTimeMap = make(map[string]int64)

// 定义接口
type ISocialDB interface {
	//更新关注表数据
	UpdateFollow(follow *model.Follow) (int64, error)

	//根据关注者查询被关注者s
	GetFolloweds(follower int64) ([]int64, error)

	//根据被关注者查询关注者s
	GetFollowers(followed int64) ([]int64, error)

	//根据关注者查询相互关注的用户(与关注者互相关注的用户)--friend
	GetFriends(follower int64) ([]int64, error)

	//向信息表中插入数据,并向redis中有序集合插入数据(键为fromuser-touser,分数为crreate_time时间戳,值为message(序列化后的json数据))
	InsertMessage(message *model.Message) error

	//根据fromuser与touser查询信息s
	GetMessages(fromUser, toUser, preMsgTime int64) ([]*social.Message, error)

	//查询/创建room((查询条件(user_id=userId&&friend_id=friendId)||(user_id=friendId&&friend_id=userId))查询room,存在返回room,不存在创建room并返回room)
	GetRoom(userId, friendId int64) (*model.Room, error)

	//创建room
	CreateRoom(userId, friendId int64) error

	//通过userId与friendIds批量查询roomids(查询条件(user_id=userId&&friend_id=friendIds)||(user_id=friendIds&&friend_id=userId))
	GetRoomIds(userId int64, friendIds []int64) ([]int64, error)

	//通过roomIds批量查询最新1条的message
	GetNewMessages(roomIds []int64) ([]social.Message, error)

	//关闭redis
	Close()

	// 批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map
	GetFollows(req *social.FindFollowsRequest) (map[int64]bool, error)

	// 弃用-巡回扫描的请求方式-难以设计消息队列的消费关系/设计出来消耗太大
	// 将message数据添加到消息队列中
	InsertMessageToQueue(message *model.Message) error

	// 从消息队列中获取message数据
	GetMessageFromQueue(roomId int64) (*model.Message, error)
}

// 定义结构体
type SocialDB struct {
	db   *gorm.DB
	pool *redis.Pool
}

// 确保结构体实现接口的小技巧:
var _ ISocialDB = (*SocialDB)(nil)

// 实例化
func NewSocialDB(sql *db.Mysql, rds *db.Redis) (socialDB *SocialDB) {
	socialDB = &SocialDB{db.NewDB(sql), db.NewRedisPool(rds)}
	return
}

func (s *SocialDB) Close() {
	s.pool.Close()
	db.CloseDB(s.db)
}

// 更新关注表数据
func (s *SocialDB) UpdateFollow(follow *model.Follow) (int64, error) {
	f := *follow
	//查询/创建关注记录
	result := s.db.FirstOrCreate(&follow, model.Follow{Follower: follow.Follower, Followed: follow.Followed})
	//是否需要更新
	if f.IsFollow != follow.IsFollow {
		result = s.db.Model(&model.Follow{}).Where(model.Follow{Follower: follow.Follower, Followed: follow.Followed}).Update("is_follow", f.IsFollow)
		return 1, result.Error
	}
	//未关注时,取消关注
	if result.RowsAffected == 1 && follow.IsFollow == 2 {
		return 0, result.Error
	}
	return result.RowsAffected, result.Error
}

// 根据关注者查询被关注者s
func (s *SocialDB) GetFolloweds(follower int64) (followeds []int64, err error) {
	err = s.db.Model(&model.Follow{}).Where("follower = ? and is_follow = ?", follower, 1).Pluck("followed", &followeds).Error
	if err != nil {
		logger.Debugf("GetFollowedsByFollower err:%v", err)
	}
	return
}

// 根据被关注者查询关注者s
func (s *SocialDB) GetFollowers(followed int64) (followers []int64, err error) {
	err = s.db.Model(&model.Follow{}).Where("followed = ? and is_follow = ?", followed, 1).Pluck("follower", &followers).Error
	if err != nil {
		logger.Debugf("GetFollowersByFollowed err:%v", err)
	}
	return
}

// 根据关注者查询相互关注的用户(与关注者互相关注的用户)--friend
func (s *SocialDB) GetFriends(follower int64) (friends []int64, err error) {
	// 获取关注者关注的用户
	followeds, err := s.GetFolloweds(follower)
	if err != nil {
		logger.Debugf("GetFollowedsByFollower err:%v", err)
		return
	}
	//根据以被关注者followeds作为关注者查询关注了关注者follower的用户
	err = s.db.Model(&model.Follow{}).Where("follower in (?) and followed = ? and is_follow = ?", followeds, follower, 1).Pluck("follower", &friends).Error
	if err != nil {
		logger.Debugf("GetFriendsByFollower err:%v", err)
		return
	}
	//返回数据
	return
}

// 查询/创建room((查询条件(user_id=userId&&friend_id=friendId)||(user_id=friendId&&friend_id=userId))查询room,存在返回room,不存在创建room并返回room)
func (s *SocialDB) GetRoom(userId, friendId int64) (room *model.Room, err error) {
	//查询条件(user_id=userId&&friend_id=friendId)||(user_id=friendId&&friend_id=userId)
	//查询room
	room = &model.Room{UserId: userId, FriendId: friendId}
	err = s.db.Model(&model.Room{}).Where("(user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)", userId, friendId, friendId, userId).FirstOrCreate(room).Error
	return
}

// 创建room
func (s *SocialDB) CreateRoom(userId, friendId int64) error {
	room := &model.Room{
		UserId:   userId,
		FriendId: friendId,
	}
	return s.db.Create(room).Error
}

// 向信息表中插入数据,并向redis中有序集合插入数据(键为fromuser-touser,分数为creat_date时间戳,值为message(序列化后的json数据))
func (s *SocialDB) InsertMessage(message *model.Message) error {
	//开启事务
	tx := s.db.Begin()
	//插入数据并获取到存入后的实例数据
	err := tx.Create(message).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf("tx.Create err:%v", err)
		return err
	}

	// 已经弃用___信息不再使用redis___________________________________________________________________________________________________________________
	// //序列化message
	// jsonData, err := jsoner.Marshal(message)
	// if err != nil {
	// 	tx.Rollback()
	// 	logger.Errorf("json.Marshal err:%v", err)
	// 	return err
	// }

	// //向redis中有序集合插入数据(键为fromuser-touser,分数为crreate_time时间戳,值为message(序列化后的json数据))
	// _, err = s.conn.Do("ZADD", fmt.Sprintf("%d-%d", message.FromUserId, message.ToUserId), message.CreateDate, jsonData)
	// if err != nil {
	// 	tx.Rollback()
	// 	logger.Errorf("conn.Do err:%v", err)
	// 	return err
	// }
	//提交事务
	tx.Commit()
	logger.Debugf("InsertMessage success")
	return nil
}

// 根据fromuser与touser查询信息s
func (s *SocialDB) GetMessages(fromUser, toUser, preMsgTime int64) (messages []*social.Message, err error) {
	//查询room
	room, err := s.GetRoom(fromUser, toUser)
	if err != nil {
		logger.Errorf("GetRoom err:%v", err)
		return
	}
	var messagesData []model.Message
	//如果preMsgTime为0,使用全局变量中的时间
	// if preMsgTime == 0 {
	// 	//获取map中的时间
	// 	preMsgTime = getMessageTimeMap[getMapString(room.RoomId, fromUser)]
	// 	logger.Debugf("获取到map中时间戳 : %d", preMsgTime)
	// }

	//查询信息
	err = s.db.Model(&model.Message{}).Where("room_id = ? and create_time > ?", room.RoomId, preMsgTime).Order("create_time asc").Limit(20).Find(&messagesData).Error
	if err != nil {
		logger.Errorf("GetMessages err:%v", err)
		return
	}
	// //更新map中的时间
	// if len(messagesData) > 0 {
	// 	getMessageTimeMap[getMapString(room.RoomId, fromUser)] = messagesData[0].Message.CreateTime
	// }

	lenth := len(messagesData)
	//如果不是第一次刷取数据,并且只有一条数据,数据的fromuser与请求的fromuser相同,则不返回数据(避免发送消息出现重复数据)
	if preMsgTime != 0 && lenth == 1 && messagesData[0].FromUserId == fromUser {
		messagesData[0].Message.Content = ""
	}
	messages = make([]*social.Message, lenth)
	for i := 0; i < lenth; i++ {
		messages[i] = messagesData[i].ToMessage()
	}

	// 已经弃用___信息不再使用redis___________________________________________________________________________________________________________________
	// //从redis中有序集合中查询数据
	// messagesJson, err := redis.ByteSlices(s.conn.Do("ZRANGEBYSCORE", fmt.Sprintf("%d-%d", fromUser, toUser), "-inf", "+inf"))
	// if err != nil {
	// 	return
	// }
	// messages = make([]social.Message, 0, len(messagesJson))
	// //反序列化数据
	// for i := 0; i < len(messagesJson); i++ {
	// 	//json数据不为空
	// 	if messagesJson[i] == nil {
	// 		//定义一个空的messagedata
	// 		var message model.Message
	// 		err = jsoner.Unmarshal(messagesJson[i], &message)
	// 		if err != nil {
	// 			return
	// 		}
	// 		messages[i] = *message.ToMessage()
	// 	}
	// }

	return
}

// 通过userId与friendIds批量查询roomids(查询条件(user_id=userId&&friend_id=friendIds)||(user_id=friendIds&&friend_id=userId))
func (s *SocialDB) GetRoomIds(userId int64, friendIds []int64) (roomIds []int64, err error) {
	//查询room
	err = s.db.Model(&model.Room{}).Where("(user_id = ? and friend_id in (?)) or (user_id in (?) and friend_id = ?)", userId, friendIds, friendIds, userId).Pluck("room_id", &roomIds).Error
	if err != nil {
		logger.Errorf("GetRoomIds err:%v", err)
		return
	}
	return
}

// 通过roomIds批量查询最新1条的message
func (s *SocialDB) GetNewMessages(roomIds []int64) (messages []social.Message, err error) {
	var messagesData []model.Message
	//查询信息
	err = s.db.Model(&model.Message{}).
		// Select("message.*").
		// Where("room_id IN (?)", roomIds).
		// Order("create_time DESC").
		Joins("JOIN (SELECT room_id, MAX(create_time) AS create_time FROM messages WHERE room_id IN (?) GROUP BY room_id) AS sub ON messages.room_id = sub.room_id AND messages.create_time = sub.create_time", roomIds).
		Find(&messagesData).Error

	if err != nil {
		logger.Errorf("GetNewMessages err:%v", err)
		return
	}
	messages = make([]social.Message, len(messagesData))
	for i := 0; i < len(messagesData); i++ {
		messages[i] = *messagesData[i].ToMessage()
	}
	return
}

// 批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map
func (s *SocialDB) GetFollows(req *social.FindFollowsRequest) (map[int64]bool, error) {
	//查询关注者
	var follows []model.Follow
	err := s.db.Model(&model.Follow{}).Where("follower = ? and is_follow = ? and followed IN (?)", req.UserId, 1, req.Ids).Find(&follows).Error
	if err != nil {
		logger.Debugf("GetFollowsByFollower err:%v", err)
		return nil, err
	}
	//定义map
	followsMap := make(map[int64]bool)
	for _, follow := range follows {
		followsMap[follow.Followed] = true
	}
	return followsMap, nil
}

// 弃用-巡回扫描的请求方式-难以设计消息队列的消费关系/设计出来消耗太大
// 将message数据添加到消息队列中
func (s *SocialDB) InsertMessageToQueue(message *model.Message) error {
	//获取链接
	conn := s.pool.Get()
	defer conn.Close()
	//序列化message
	jsonData, err := jsoner.Marshal(message)
	if err != nil {
		logger.Errorf("json.Marshal err:%v", err)
		return err
	}
	//将message数据添加到消息队列中
	_, err = conn.Do("LPUSH", GetQueueKey(message.RoomId), jsonData)
	if err != nil {
		logger.Errorf("conn.Do err:%v", err)
		return err
	}
	return nil
}

// 从消息队列中获取message数据
func (s *SocialDB) GetMessageFromQueue(roomId int64) (message *model.Message, err error) {
	//获取链接
	conn := s.pool.Get()
	defer conn.Close()
	//从消息队列中获取message数据
	messageJson, err := redis.Bytes(conn.Do("RPOP", GetQueueKey(roomId)))
	if err != nil {
		logger.Errorf("redis.Bytes err:%v", err)
		return
	}
	//反序列化message
	err = jsoner.Unmarshal(messageJson, &message)
	if err != nil {
		logger.Errorf("json.Unmarshal err:%v", err)
		return
	}
	return
}

// 拼接消息队列房间id
func GetQueueKey(roomId int64) string {
	return fmt.Sprintf("message_queue:%d", roomId)
}

// 拼接时间map字符串----------------------------------弃用
func getMapString(roomId int64, userId int64) string {
	return fmt.Sprintf("%d:%d", roomId, userId)
}

//聊天功能;
