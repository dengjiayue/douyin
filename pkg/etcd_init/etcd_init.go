package etcdinit

import (
	"douyin/pkg/etcd_register"
	"douyin/pkg/logger"
)

// 初始化etcd
func InitETCD(path, addr string, cheektime int64) *etcd_register.EtcdRegister {
	// --------------------------------------------------------------------------------------------

	//初始化etcd
	etcdRegister, err := etcd_register.NewEtcdRegister()
	if err != nil {
		logger.Errorf("初始化etcd失败: %s\n", err)
		panic(err)
	}

	//注册服务
	err = etcdRegister.RegisterServer(path, addr, cheektime)
	if err != nil {
		logger.Errorf("注册服务失败: %s\n", err)
		panic(err)
	}
	return etcdRegister
}
