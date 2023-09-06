package db

import (
	"testing"
)

// 定义User模型，绑定users表，ORM库操作数据库，需要定义一个struct类型和MYSQL表进行绑定或者叫映射，struct字段和MYSQL表字段一一对应
// 在这里User类型可以代表mysql users表
// type User struct {
// 	ID int // 主键
// 	// 通过在字段后面的标签说明，定义golang字段和表字段的关系
// 	// 例如 `gorm:"column:username"` 标签说明含义是: Mysql表的列名（字段名)为username
// 	// 这里golang定义的Username变量和MYSQL表字段username一样，他们的名字可以不一样。
// 	Username string `gorm:"column:username"`
// 	email    string `gorm:"column:email"`
// 	// 创建时间，时间戳
// 	age int `gorm:"column:age"`
// }

// // 设置表名，可以通过给struct类型定义 TableName函数，返回当前struct绑定的mysql表名是什么
// func (u User) TableName() string {
// 	// 绑定MYSQL表名为users
// 	return "users"
// }

func TestNewRedis(t *testing.T) {

	// // 测试redis带密码连接
	// conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	// if err != nil {
	// 	fmt.Println("redis连接失败")
	// 	panic(err)
	// }
	// // 身份验证密码
	// defer conn.Close()
	// if _, authErr := conn.Do("AUTH", "douyinapp"); authErr != nil {
	// 	fmt.Println("Failed to authenticate:", authErr)
	// 	return
	// }

	// //简单的set和get操作
	// _, err = conn.Do("set", "name", "dengjiayue")
	// if err != nil {
	// 	fmt.Println("redis set failed:", err)
	// } else {
	// 	fmt.Println("redis set success")
	// }
	// name, err := redis.String(conn.Do("get", "name"))
	// if err != nil {
	// 	fmt.Println("redis get failed:", err)
	// } else {
	// 	fmt.Printf("get name: %v\n", name)
	// }
	// //删除riedis中的数据
	// conn.Do("del", "name")
	// rds := Redis{
	// 	Host:     "127.0.0.1",
	// 	Password: "",
	// 	Port:     6379,
	// }
	// // 连接redis
	// conn := NewRedis(&rds)
	// defer conn.Close()
	// //从redis中查询
	// ids := []int64{1, 2, 5, 6}
	// fmt.Printf("data=%#v\n", ids[:len(ids)])
	// data := []interface{}{"test", ids[0], "MAKE", ids[1], "sjy", ids[2], "zs", ids[3], "ls"}
	// //批量插入数据
	// conn.Do("HMSET", data...)

	// userIds := []interface{}{0, ids[0], ids[1], ids[2], ids[3], ids[0]}
	// args := []interface{}{"test"}
	// args = append(args, userIds...)
	// //检查args
	// fmt.Println(args)
	// userJSON, err := redis.ByteSlices(conn.Do("HMGET", args...))
	// if err != nil {
	// 	fmt.Println("redis get failed:", err)
	// } else {
	// 	if userJSON[0] != nil {
	// 		fmt.Printf("jsondata=%v|||\n", userJSON[0])
	// 	}

	// 	// //解析json
	// 	// var users model.UserData
	// 	// err := jsoner.Unmarshal(userJSON[0], &users)
	// 	// if err != nil {
	// 	// 	fmt.Println("json unmarshal failed:", err)
	// 	// 	return
	// 	// }

	// 	// fmt.Printf("user=%#v\n", users)
	// }
	// conn.Do("HDEL", "test")

	// conn.Do("HDEL", "users")

	//测试批量插入数据

	sql := Mysql{"root", "my password", "我的地址", 3306, "video"}
	db := NewDB(&sql)
	// 关闭数据库连接
	defer CloseDB(db)
	//清除video_list表的内容
	db.Exec("truncate table video_list")
	// // 通过gorm操作数据库
	// user := &User{}
	// db.Where("id=?", 1).Find(user)
	// fmt.Printf("u=ser=%#v\n", user)
	// // 通过Do函数向redis写入数据
	// conn.Do("set", "name", "dengjiayue")
	// name, err := redis.String(conn.Do("get", "name"))
	// if err != nil {
	// 	fmt.Println("redis get failed:", err)
	// } else {
	// 	fmt.Printf("get name: %v\n", name)
	// }
	// //删除riedis中的数据
	// conn.Do("del", "name")
	// // 再次读取
	// name, err = redis.String(conn.Do("get", "name"))
	// if err != nil {
	// 	fmt.Println("redis get failed:", err)
	// } else {
	// 	fmt.Printf("get name: %v\n", name)
	// }
}
