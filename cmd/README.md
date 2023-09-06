编译执行文件
```cmd
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /root/gop1/douyinapp/cmd/gateway ./cmd/gateway/main.go
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /root/gop1/douyinapp/cmd/user ./cmd/user/main.go
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /root/gop1/douyinapp/cmd/video ./cmd/video/main.go
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /root/gop1/douyinapp/cmd/video_list ./cmd/video_list/main.go
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /root/gop1/douyinapp/cmd/social ./cmd/social/main.go
```


```cmd
scp -r /root/gop1/douyinapp/cmd djy@djy.bmft.tech:/home/djy/douyinapp
```

* 启动 etcd

拉取镜像
```bash
docker pull quay.io/coreos/etcd
```
* 启动 etcd(无身份验证)
```bash
docker run -d -p 2379:2379 -p 2380:2380 --name etcd-server quay.io/coreos/etcd /usr/local/bin/etcd --name my-etcd-1 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 --initial-advertise-peer-urls http://0.0.0.0:2380 --listen-peer-urls http://0.0.0.0:2380 --initial-cluster my-etcd-1=http://0.0.0.0:2380
```

* 增加etcd身份验证

```bash
docker run -v /root/gop1/douyin/pkg/etcd_client/etcd.conf:/etc/etcd/etcd.conf -d -p 2379:2379 -p 2380:2380 --name etcd-server quay.io/coreos/etcd /usr/local/bin/etcd --name my-etcd-1 --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 --initial-advertise-peer-urls http://0.0.0.0:2380 --listen-peer-urls http://0.0.0.0:2380 --initial-cluster my-etcd-1=http://0.0.0.0:2380
```

其中/root/gop1/douyin/pkg/etcd_client/etcd.conf 是本地的配置文件路径, 里面配置了 etcd 的集群信息

```conf
security:
  authentication: "simple"
users:
  - name: "dengjiayue"
    password: "douyinapp"

```

* 配置etcd的用户密码

```
docker exec -it etcd-server /bin/sh
#         
# etcdctl user add dengjiayue:douyinapp
User dengjiayue created
# etcdctl role add master
Role master created
# etcdctl role grant master --path=/ --readwrite
Role master updated
# etcdctl user grant dengjiayue --roles=maste
Role master is granted to user dengjiayue
```

重启etcd

```bash
docker restart etcd-server
```

* 启动redis

运行redis镜像(无密码)
```bash
docker run -d --name my-redis -p 6379:6379 redis
```

运行redis镜像(有密码)
```bash
docker run -d --name my-redis -p 6379:6379 -e REDIS_PASSWORD=douyinapp redis:latest --requirepass douyinapp
```

docker重启
```bash
docker restart <服务名/id>
```

docker停止
```bash
docker stop <服务名/id>
```

docker删除服务
```bash
docker rm <服务名/id>
```

