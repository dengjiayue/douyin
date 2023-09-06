

* 更改video_list的go文件
将video_list.go中video类的VideoAuthor字段
与video.go中的comment类的Author字段
(解释:不改会引起gorm的查询id冲突报错)
```go
//video_list.go
	VideoAuthor   *user.User `protobuf:"bytes,2,opt,name=video_author,json=videoAuthor,proto3" json:"video_author,omitempty" gorm:"foreignKey:Id"` // 使用 foreignKey 标签指定关联关系

//comment.go
	User       *user.User `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty" gorm:"foreignKey:Id"`                               // 评论用户信息

```

* 导包错误修复
```go
	user "douyin/api/pb/user"
	video_list "douyin/api/pb/video_list"
```

```bash
protoc -I ./api --go_out=./api --go-grpc_out=./api api/proto/*.proto
```
