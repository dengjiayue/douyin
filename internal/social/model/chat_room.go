package model

// 定义结构体
type Room struct {
	RoomId   int64 `gorm:"column:room_id;type:bigint(20) unsigned;not null;primary_key;AUTO_INCREMENT"`
	UserId   int64 `gorm:"column:user_id;type:bigint(20) unsigned;not null"`
	FriendId int64 `gorm:"column:friend_id;type:bigint(20) unsigned;not null"`
}

// 定义表名
func (Room) TableName() string {
	return "chat_room"
}
