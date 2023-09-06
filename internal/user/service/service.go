package service

import (
	"douyin/internal/user/db"
)

type IUserService interface {
	// 继承db.IUserDB
	db.IUserDB
}

// type UserService struct {
// 	UserDB *db.UserDB
// }

func NewUserService(userDB db.IUserDB) IUserService {
	return userDB
}
