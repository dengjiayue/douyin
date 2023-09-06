package handler

import (
	"douyin/api/pb/user"
	"douyin/internal/gateway/service"
	my_jwt "douyin/pkg/jwt"
	"douyin/pkg/logger"
	"douyin/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin
func UserGin(s *web.Web, uh *UserHandler) {
	//ping测试
	s.GET("/ping", uh.Ping)

	logger.Debugf("注册 UserServer")
	s.POST("/douyin/user/register/", uh.DouyinUserRegister)

	logger.Debugf("登录 UserServer")
	s.POST("/douyin/user/login/", uh.DouyinUserLogin)

	logger.Debugf("获取用户信息 UserServer")
	s.GET("/douyin/user/", uh.DouyinUserInfo)
}

type UserHandler struct {
	user.UnimplementedUserServer
	//
	UserClient service.IUserClient
}

func NewUserHandler(UserClient *service.IUserClient) *UserHandler {
	return &UserHandler{
		UserClient: *UserClient,
	}
}

func (s *UserHandler) DouyinUserRegister(c *gin.Context) {
	var req user.DouyinUserRegisterRequest
	//username与password长度校验,在32位以内
	if err := c.ShouldBindQuery(&req); err != nil || len(req.Username) > 32 || len(req.Password) > 32 {
		logger.Errorf("[error] 参数异常： %s\n", err)
		c.JSON(400,
			&user.DouyinUserResponse{
				StatusCode: 400,
				StatusMsg:  "参数异常",
				User:       nil,
			})
		return
	}

	resp, err := s.UserClient.RegisterUser(&req)
	if err != nil {
		logger.Errorf("[error] 调用用户服务异常： %s\n", err)
		c.JSON(500,
			resp)
		return
	}
	//注册成功,生成token
	token, err := my_jwt.GenToken(resp.UserId)
	if err != nil {
		logger.Errorf("[error] 生成token异常: %s\n", err)
		c.JSON(500,
			&user.DouyinUserResponse{
				StatusCode: 500,
				StatusMsg:  "生成token异常",
				User:       nil,
			})
		return
	} else if resp.StatusCode != 0 {
		c.JSON(500, resp)
		return
	}
	//将token储存到redis中
	if !s.UserClient.SetLogInStatus(&resp.UserId, &token) {
		logger.Errorf("[error] token储存redis失败: %s\n", err)
		c.JSON(500, &user.DouyinUserResponse{
			StatusCode: 500,
			StatusMsg:  "注册成功,但token储存redis失败",
			User:       nil,
		})
		return
	}
	resp.StatusCode = 0
	resp.Token = token
	c.JSON(http.StatusOK, resp)
}

func (s *UserHandler) DouyinUserLogin(c *gin.Context) {
	var req user.DouyinUserLoginRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("[error] 参数异常： %s\n", err)
		c.JSON(400, &user.DouyinUserResponse{
			StatusCode: 400,
			StatusMsg:  "参数异常",
			User:       nil,
		})
		return
	}

	resp, err := s.UserClient.LoginUser(&req)
	if err != nil {
		logger.Errorf("[error] 调用用户服务异常： %s\n", err)
		c.JSON(500, &user.DouyinUserResponse{
			StatusCode: 500,
			StatusMsg:  "调用用户服务异常",
			User:       nil,
		})
		return
	}

	//登录成功,生成token
	token, err := my_jwt.GenToken(resp.UserId)
	if err != nil {
		logger.Errorf("[error] 生成token异常: %s\n", err)
		c.JSON(500, &user.DouyinUserResponse{
			StatusCode: 500,
			StatusMsg:  "生成token异常",
			User:       nil,
		})
		return
	} else if resp.StatusCode != 0 {
		c.JSON(400, resp)
		return
	} else {
		resp.Token = token
		//将token储存到redis中
		if !s.UserClient.SetLogInStatus(&resp.UserId, &token) {
			logger.Errorf("[error] token储存redis失败")
			c.JSON(500, &user.DouyinUserResponse{
				StatusCode: 500,
				StatusMsg:  "登录失败",
				User:       nil,
			})
			return
		}
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 查询用户信息
func (s *UserHandler) DouyinUserInfo(c *gin.Context) {
	var req user.DouyinUserRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Errorf("[error] 参数异常： %s\n", err)
		c.JSON(400, &user.DouyinUserResponse{
			StatusCode: 400,
			StatusMsg:  "参数异常",
			User:       nil,
		})
		return
	}

	uid, b := c.Get("user_id")
	id := uid.(int64)
	if id == 0 || !b {
		logger.Errorf("[error] 用户鉴权失败")
		c.JSON(400, &user.DouyinUserResponse{
			StatusCode: 400,
			StatusMsg:  "用户鉴权失败",
			User:       nil,
		})
		return
	}
	req.Token = ""
	resp, err := s.UserClient.UserInfo(&req)
	if err != nil {
		logger.Errorf("[error] 调用用户服务异常： %s\n", err)
		c.JSON(500, &user.DouyinUserResponse{
			StatusCode: 500,
			StatusMsg:  "调用用户服务异常",
			User:       nil,
		})
		return
	}
	resp.StatusCode = 0
	c.JSON(http.StatusOK, resp)
}

// 实现ping测试
func (s *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
