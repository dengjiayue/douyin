
```log
[2023-08-07 00:35:11.331]       DEBUG   db/idb.go:62    数据同步到redis完成
[2023-08-07 00:35:11.331]       DEBUG   handler/handler.go:13   注册 VideoServer
2023/08/07 00:35:11 /root/gop1/douyin/internal/video_list/db/idb.go:48
[error] invalid field found for struct douyin/internal/video_list/model.VideoData's field VideoAuthor: define a valid foreign key for relations or implement the Valuer/Scanner interface
[2023-08-07 00:35:11.332]       DEBUG   db/idb.go:62    数据同步到redis完成
[2023-08-07 00:35:11.333]       DEBUG   handler/handler.go:13   注册 VideoServer
2023/08/07 00:35:11 bindLease success &{cluster_id:11588568905070377092 member_id:128088275939295631 revision:84 raft_term:4  <nil> {} [] 0} 
[2023-08-07 00:35:15.918]       DEBUG   db/idb.go:73    收到请求:0
[2023-08-07 00:35:15.919]       DEBUG   db/idb.go:116   检查视频信息:[]model.VideoData(nil)
[2023-08-07 00:35:15.919]       ERROR   db/idb.go:120   rpc调用失败:rpc error: code = Unavailable desc = last connection error: connection error: desc = "transport: Error while dialing: dial tcp: address tcp////: unknown port"

2023/08/07 09:53:38 /root/gop1/douyin/internal/video_list/db/idb.go:48 Error 1146 (42S02): Table 'video.video' doesn't exist
[8.981ms] [rows:0] SELECT * FROM `video` WHERE create_time<0 ORDER BY create_time desc


2023/08/07 09:59:22 /root/gop1/douyin/internal/video_list/db/idb.go:48
[error] invalid field found for struct douyin/internal/video_list/model.VideoData's field VideoAuthor: define a valid foreign key for relations or implement the Valuer/Scanner interface

```

