package model

import "douyin/api/pb/user"

type UserData struct {
	*user.User
	Password string `json:"password" gorm:"not null"`
}

// 指定表名
func (UserData) TableName() string {
	return "users"
}

type UserDB struct {
}
