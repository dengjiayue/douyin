package discover

import (
	"context"
	"douyin/pkg/logger"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type etcdResolverBuilder struct {
	etcdClient *clientv3.Client
}

func NewEtcdResolverBuilder() *etcdResolverBuilder {

	// 创建etcd客户端连接
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
		Password:    "douyinapp",
		Username:    "dengjiayue",
	})

	if err != nil {
		logger.Errorf("连接异常： %s\n", err)
		panic(err)
	}

	return &etcdResolverBuilder{
		etcdClient: etcdClient,
	}
}

func (erb *etcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	// 获取指定前缀的etcd节点值
	// /ns->/ns/order-service-1   /ns/order-service-2
	prefix := target.URL.Path

	logger.Infof("prefix: %s", prefix)

	// 获取 etcd 中服务保存的ip列表
	res, err := erb.etcdClient.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		logger.Errorf("获取异常： %s\n", err)
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	es := &etcdResolver{
		cc:         cc,
		etcdClient: erb.etcdClient,
		ctx:        ctx,
		cancel:     cancelFunc,
		scheme:     target.URL.Path,
	}

	// 将获取到的ip和port保存到本地的map中
	logger.Infof("获取到的服务列表：%v", res.Kvs)
	for _, kv := range res.Kvs {
		es.store(kv.Key, kv.Value)
	}

	// 更新拨号里的ip列表
	es.updateState()

	// 监听etcd中的服务是否变化

	go es.watcher()
	return es, nil
}

func (erb *etcdResolverBuilder) Scheme() string {
	return "etcd"
}
