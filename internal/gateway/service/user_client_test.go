package service

//-----------------------------------------------------弃用测试

// func TestUserClient_RegisterUser(t *testing.T) {
// 	config.Init("../../../configs/gateway.yaml")
// 	logger.Init(config.GlobalConfig.Log)

// 	conn, err := grpc_client.NewGrpcClient("localhost:10001")
// 	if err != nil {
// 		t.Errorf("err: %s", err)
// 	}
// 	redisConn := db.NewRedisPool(&db.Redis{}).Get()
// 	userClient := NewUserClient(redisConn, conn)
// 	resp, err := userClient.RegisterUser(&user.DouyinUserRegisterRequest{
// 		Username: "test",
// 		Password: "test",
// 	})
// 	if err != nil {
// 		t.Errorf("err: %s", err)
// 	}

// 	t.Logf("resp: %v", resp)
// }
