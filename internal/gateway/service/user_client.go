package service

import (
	"context"
	"crypto/sha256"
	"douyin/api/pb/user"
	"douyin/internal/gateway/db"
	"douyin/pkg/grpc_client"
	"encoding/hex"

	"github.com/gomodule/redigo/redis"
)

// IUserClient Grpc 调用用户服务
type IUserClient interface {
	//继承db的接口
	db.IGetWay
	// 注册用户
	RegisterUser(req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error)
	// 登录用户
	LoginUser(req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error)
	// 查询用户信息
	UserInfo(req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error)
}

type UserClient struct {
	db.IGetWay
	UserClient user.UserClient
}

func NewUserClient(pool *redis.Pool, conn *grpc_client.GrpcClient) IUserClient {

	return &UserClient{
		UserClient: user.NewUserClient(conn.ClientConn),
		IGetWay:    &db.GetWay{Pool: pool},
	}
}

// 密码哈希加密
func encryptPassword(password string) string {
	// 将字符串类型的密码转换为字节类型
	passwordBytes := []byte(password)
	// 创建SHA-256哈希对象
	sha256Hash := sha256.New()
	// 将密码字节数据传入哈希对象
	sha256Hash.Write(passwordBytes)
	// 获取哈希值的字节数据
	hashBytes := sha256Hash.Sum(nil)
	// 将字节数据转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)
	// 返回十六进制字符串类型的哈希值
	return hashString
}

// 注册用户
func (s *UserClient) RegisterUser(req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	// 加密密码
	req.Password = encryptPassword(req.Password)
	resp, err = s.UserClient.DouyinUserRegister(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 登录用户
func (s *UserClient) LoginUser(req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	// 加密密码
	req.Password = encryptPassword(req.Password)
	resp, err = s.UserClient.DouyinUserLogin(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 查询用户信息
func (s *UserClient) UserInfo(req *user.DouyinUserRequest) (resp *user.DouyinUserResponse, err error) {

	resp, err = s.UserClient.DouyinUserInfo(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
