package db

import (
	"testing"
)

func TestNewRedis(t *testing.T) {

	//测试批量插入数据

	sql := Mysql{"root", "my password", "我的地址", 3306, "video"}
	db := NewDB(&sql)
	// 关闭数据库连接
	defer CloseDB(db)
	//清除video_list表的内容
	db.Exec("truncate table video_list")

}
