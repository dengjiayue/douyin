package model

import (
	"douyin/api/pb/video"
	"time"
)

// 定义comment表结构体,用于查询(create_time为时间戳 ,默认值为当前时间的时间戳)
type CommentData struct {
	//主键id:
	Id int64 `json:"id" gorm:"primary_key;not null"`
	//视频id
	VideoId int64 `json:"video_id" gorm:"not null"`
	//评论人id
	UserId int64 `json:"user_id" gorm:"not null"`
	//评论内容
	Content string `json:"content" gorm:"not null"`
	//评论时间: 默认值为当前时间的时间戳
	CreateTime int64 `json:"create_time" gorm:"not null;default:0"`
}

// 指定表名
func (CommentData) TableName() string {
	return "comments"
}

// CommentData转换为Comment
func (commentData *CommentData) ToComment() (comment *video.Comment) {
	comment = &video.Comment{
		Id:         commentData.Id,
		Content:    commentData.Content,
		CreateDate: commentData.TimeStampToString(),
	}
	return
}

// 时间戳转换为时间字符串(mm-dd)
func (commentData *CommentData) TimeStampToString() (Time string) {
	//时间戳转换为时间字符串
	// 将时间戳转换为时间
	t := time.Unix(commentData.CreateTime, 0)

	// 格式化为 mm-dd
	Time = t.Format("01-02")

	return
}
