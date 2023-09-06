package service

import "douyin/internal/video_list/db"

// 继承db.IVideoDB
type IVideoService interface {
	db.IVideoListDB
}

// 实例化
func NewVideoListService(videoDB db.IVideoListDB) IVideoService {
	return videoDB
}
