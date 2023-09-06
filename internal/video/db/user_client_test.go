package db

//已经将客户端包放入到了internal/user/db/user_client.go中,所以这里不需要再定义客户端

// func TestNewUserClient(t *testing.T) {
// 	//初始化日志,避免空指针报错
// 	logger.Init(nil)
// 	client := NewUserClient()
// 	data, err := client.UserInfo(&user.DouyinUserRequest{UserId: 5})
// 	if err != nil {
// 		fmt.Printf("err=%s\n", err)
// 		return
// 	}
// 	fmt.Printf("data=%#v\n", data.User)
// }
