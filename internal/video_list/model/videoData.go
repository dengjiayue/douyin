package model

import (
	"douyin/api/pb/video_list"
)

// VideoData 视频数据
type VideoData struct {
	*video_list.Video
	//作者id
	User_id    int64 `json:"user_id" gorm:"not null"`
	CreateTime int64 `json:"create_time"`
}

// 指定表名
func (VideoData) TableName() string {
	return "videos"
}
