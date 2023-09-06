package client

import (
	"douyin/api/pb/user"
	"douyin/pkg/logger"
	"fmt"

	"testing"
)

func TestUserClient_UsersInfo(t *testing.T) {
	logger.Init(nil)
	// 初始化
	userClient := NewUserClient()
	resp, err := userClient.UsersInfo(&user.DouyinUsersRequest{UserIds: []int64{1, 5, 6}})
	if err != nil {
		fmt.Printf("err: =%s\n", err)
		return
	}
	fmt.Printf("resp: =%#v\n", resp)
}
