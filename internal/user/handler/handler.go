package handler

import (
	"context"
	"douyin/api/pb/user"
	"douyin/internal/user/service"
	"douyin/pkg/grpc_server"
	"douyin/pkg/logger"
)

func Register(s *grpc_server.GrpcServer, uh *UserHandler) {
	logger.Debugf("注册 UserServer")
	user.RegisterUserServer(s, uh)
}

type UserHandler struct {
	user.UnimplementedUserServer

	UserService service.IUserService
}

func NewUserHandler(UserService *service.IUserService) *UserHandler {
	return &UserHandler{
		UserService: *UserService,
	}
}

func (s *UserHandler) DouyinUserRegister(ctx context.Context, req *user.DouyinUserRegisterRequest) (*user.DouyinUserRegisterResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.UserService.RegisterUser(req)
}

func (s *UserHandler) DouyinUserLogin(ctx context.Context, req *user.DouyinUserLoginRequest) (*user.DouyinUserLoginResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.UserService.LoginUser(req)
}

// 查询用户信息
func (s *UserHandler) DouyinUserInfo(ctx context.Context, req *user.DouyinUserRequest) (*user.DouyinUserResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.UserService.UserInfo(req)
}

// 批量查询用户信息
func (s *UserHandler) DouyinUsersInfo(ctx context.Context, req *user.DouyinUsersRequest) (*user.DouyinUsersResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	return s.UserService.UsersInfo(req)
}

// 更新用户数据
func (s *UserHandler) DouyinUserChange(ctx context.Context, req *user.DouyinUserChangeRequest) (*user.DouyinUserChangeResponse, error) {
	// 参数校验
	logger.Debugf("req: %+v\n", req)
	err := s.UserService.UserChange(req)
	if err != nil {
		return nil, err
	}
	return &user.DouyinUserChangeResponse{StatusCode: 200}, nil
}
