package db

import (
	my_jwt "douyin/pkg/jwt"
	"douyin/pkg/logger"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// 定义接口:
type IGetWay interface {
	GetLogInStatus(id *int64, token *string) bool
	SetLogInStatus(id *int64, token *string) bool
	LogInStatus(token string) int64
	Close()
}

// 定义结构体:(redis连接)
type GetWay struct {
	Pool *redis.Pool
}

// 判断实现接口
var _ IGetWay = &GetWay{}

// 拼接字符串
func LogInStatusString(id *int64) (r string) {
	r = fmt.Sprintf("loginstatus:%d", *id)
	logger.Debugf("鉴权: %s", r)
	return
}

// 验证登录状态:
func (s *GetWay) GetLogInStatus(id *int64, token *string) bool {
	// 获取连接
	Conn := s.Pool.Get()
	//放回连接
	defer Conn.Close()
	//校验数据
	logindata := LogInStatusString(id)
	logger.Debugf("鉴权: %s", logindata)
	if len(logindata) == 0 {
		logger.Errorf("获取登录状态失败")
		return false
	}
	t, err := redis.String(Conn.Do("GET", logindata))
	if err != nil {
		logger.Errorf("获取登录状态失败: %s", err)
		return false
	}
	return *token == t
}

// 更新登录状态
func (s *GetWay) SetLogInStatus(id *int64, token *string) bool {
	// 获取连接
	Conn := s.Pool.Get()
	//放回连接
	defer Conn.Close()
	_, err := Conn.Do("SET", LogInStatusString(id), *token)
	if err != nil {
		logger.Errorf("更新登录状态失败: %s", err)
		return false
	}
	return true
}

// 鉴权函数
func (s *GetWay) LogInStatus(token string) int64 {
	// 获取连接
	Conn := s.Pool.Get()
	//放回连接
	defer Conn.Close()
	logger.Debugf("鉴权: %s", token)
	claims, err := my_jwt.ParseToken(token)
	if err != nil {
		logger.Errorf("鉴权失败")
		return 0
	}
	//验证登录状态
	if !s.GetLogInStatus(&claims.UserID, &token) {
		logger.Errorf("鉴权失败")
		return 0
	}
	return claims.UserID
}

// 关闭redis连接
func (s *GetWay) Close() {
	// 关闭redis连接池
	s.Pool.Close()
}
