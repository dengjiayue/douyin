package model

// 定义关注表结构体
type Follow struct {

	//关注者
	Follower int64 `json:"follower" gorm:"not null"`
	//被关注者
	Followed int64 `json:"followed" gorm:"not null"`

	//是否关注
	IsFollow int32 `json:"is_follow" gorm:"not null"`
}

// 指定表名
func (Follow) TableName() string {
	return "follows"
}
