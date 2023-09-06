package model

// 定义结构体:点赞数据表
type Favorite struct {
	//视频id
	VideoId int64 `gorm:"column:video_id;type:bigint(20);not null" json:"video_id"`
	//用户id
	UserId int64 `gorm:"column:user_id;type:bigint(20);not null" json:"user_id"`
	//是否点赞1点赞,2取消点赞
	IsFavorite int32 `gorm:"column:is_favorite;type:int(11);not null" json:"is_favorite"`
}

// 指定表名
func (Favorite) TableName() string {
	return "favorites"
}
