package etcd_register

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdRegister_CreateLease(t *testing.T) {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 30 * time.Second,
		Username:    "dengjiayue", // 您添加的用户名
		Password:    "douyinapp",  // 您添加的密码
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println("Failed to create etcd client:", err)
		return
	}
	defer client.Close()

	// 执行一些操作...
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = client.Put(ctx, "key", "value")
	cancel()
	if err != nil {
		fmt.Println("Failed to put key:", err)
		return
	}

	fmt.Println("Connected to etcd and performed an operation!")
}
