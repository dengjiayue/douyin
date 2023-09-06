### 重设计聊天功能(项目原接口)
1. sql设计:

* chat_room表(聊天窗表:room_id,user_id,friend_id,)
```sql
CREATE TABLE chat_room (
    room_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '聊天窗口ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    friend_id BIGINT NOT NULL COMMENT '好友ID',
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天窗口表';

```


* message表(聊天窗消息表:id,room_id,user_id,message,create_time)
```sql
CREATE TABLE message (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '消息ID',
    room_id BIGINT NOT NULL COMMENT '聊天窗口ID',
    to_user_id BIGINT NOT NULL COMMENT '用户ID',
    from_user_id BIGINT NOT NULL COMMENT '好友ID',
    content VARCHAR(255) NOT NULL COMMENT '消息内容',
    create_date BIGINT NOT NULL COMMENT '消息创建时间',
    INDEX idx_create_date_desc (create_date DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

```

2. 接口设计:
发送聊天信息:
1. 通过user_id,friend_id查询room_id,如果不存在则创建room
2. 将数据message通过room_id,user_id插入message表,创建时间为当前时间的时间戳

数据流动:网关(用户鉴权,拿到userid)->社交服务(chat_room表,获取room_id;message表,插入message信息)->网关(返回结果)

获取聊天信息:
1. 通过user_id,friend_id查询room_id,如果不存在返回空
2. 通过room_id查询message表,按照时间倒序排列

数据流动:网关(用户鉴权,拿到userid)->社交服务(chat_room表,获取room_id;message表,获取message信息列表)->网关(返回结果)



### 重设计聊天功能(使用NATS消息队列+主动推送消息)
