package model

import (
	"douyin/api/pb/social"
	"douyin/api/pb/user"
)

func UserToFriend(user *user.User) (friend *social.FriendUser) {
	friend = &social.FriendUser{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        true,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	return
}
