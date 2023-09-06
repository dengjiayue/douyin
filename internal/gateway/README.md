### 搭建整体架构

1. 网关:
* token颁发,用户鉴权,登录状态验证(中间件)
* 请求转发(rpc转发到服务端处理)


2. 优化升级:
* 聊天:
    目前:查询数据库,获取聊天记录(问题:app字段缺失,查询信息将无限刷屏)
    优化:使用redis消息队列
    消息发送:将消息存入redis消息队列,并将消息存入数据库
    无请求参数:查询redis消息队列,获取消息
    有请求参数:查询MySQL数据库,获取消息


优化:
其他服务->user(查询user数据)
user->social(查询关注数据)

```cmd
git init
git add .
git commit -m "commit"
git branch -M master
git remote add origin git@github.com:dengjiayue/douyin.git
git push -u origin master
```