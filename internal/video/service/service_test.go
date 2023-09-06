package service

import (
	"douyin/api/pb/video"
	vdb "douyin/internal/video/db"
	db "douyin/pkg/db"
	"douyin/pkg/logger"
	"fmt"
	"testing"
)

func TestNewService(t *testing.T) {
	// 初始化日志
	logger.Init(nil)
	sql := &db.Mysql{"root", "my password", "my host", 3306, "video"}
	rds := &db.Redis{
		Host:     "127.0.0.1",
		Password: "",
		Port:     6379,
	}
	s := NewService(vdb.NewVideoDB(sql, rds))
	if s == nil {
		fmt.Printf("初始化错误")
	}

	pvideo, err :=
		//获取视频发布列表
		// s.PublishList(&video.DouyinPublishListRequest{UserId: 5})   //------------测试通过
		//点赞视频favorite/action
		s.FavoriteAction(&video.RpcFavoriteActionRequest{UserId: 5, VideoId: 1, ActionType: 2}) //------------测试通过
		//获取视频点赞列表favorite/list
		// s.FavoriteList(&video.RpcFavoriteListRequest{UserId: 5}) //------------测试通过
		//评论操作comment/action
		// s.CommentAction(&video.RpcCommentActionRequest{UserId: 5, VideoId: 1, CommentText: "测试评论", ActionType: 1}) //------------测试通过
		//获取视频评论列表comment/list
		// s.CommentList(&video.RpcCommentListRequest{VideoId: 1}) //------------测试通过

	if err != nil {
		fmt.Printf("测试失败%s", err)
		t.Logf("测试失败%s", err)
		return
	}
	fmt.Printf("测试成功:%v", pvideo)
	t.Logf("测试成功:%v", pvideo)

}
