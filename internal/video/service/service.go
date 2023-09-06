package service

import "douyin/internal/video/db"

// 定义IVideoService接口
type IVideoService interface {
	//继承db.IVideoDB接口
	db.IVideoDB
}

// 实例化
func NewService(videoDB db.IVideoDB) IVideoService {
	return videoDB
}
