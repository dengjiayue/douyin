package db

import (
	"douyin/pkg/db"
	"douyin/pkg/logger"
	"testing"
)

func TestVideoListDB_Feed(t *testing.T) {
	//初始化日志
	logger.Init(nil)
	sql := &db.Mysql{Username: "root", Password: "my password", Host: "my host", Port: 3306, Dbname: "video"}
	rds := &db.Redis{
		Host:     "127.0.0.1",
		Password: "",
		Port:     6379,
	}
	videoDb := NewVideoListDB(sql, rds)
	videoDb.Sync(0)
	videoDb.Close()
	// data, err := videoDb.Feed(&video_list.DouyinFeedRequest{})
	// if err != nil {
	// 	fmt.Printf("err=%v\n", err)
	// }
	// fmt.Printf("data=%#v\n", data)
}
