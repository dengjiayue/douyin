package model

import (
	"douyin/api/pb/social"
	"time"
)

// 定义message表结构体
type Message struct {
	//继承social的message
	social.Message
	//room_id为聊天室id
	RoomId int64 `json:"room_id" gorm:"not null;default:0"`
	//弃用
	//create_date为时间戳 ,默认值为当前时间的时间戳
	// CreateDate int64 `json:"create_date" gorm:"not null;default:0"`
}

// 指定表名
func (Message) TableName() string {
	return "messages"
}

// Message转换为social.Message
func (message *Message) ToMessage() (messageData *social.Message) {
	return &message.Message
}

// ---弃用---
// 时间戳转换为时间字符串(mm-dd hh:mm)
func (message *Message) TimeStampToString() (Time string) {
	//时间戳(毫秒)转换为时间字符串
	// 将时间戳转换为时间
	t := time.UnixMilli(message.Message.CreateTime)

	// 格式化为 mm-dd hh:mm
	Time = t.Format("01-02 15:04:05")

	return
}
