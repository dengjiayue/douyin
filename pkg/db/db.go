package db

import (
	"douyin/pkg/logger"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

func NewDB(sql *Mysql) *gorm.DB {
	//查看仓库名
	fmt.Printf("仓库:%s\n", sql.Dbname)
	// 拼接数据库连接
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", sql.Username, sql.Password, sql.Host, sql.Port, sql.Dbname)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		logger.Errorf("数据库连接失败")
		panic(err)
	}
	logger.Debugf("数据库连接成功%s", dns)
	return db
}

// 关闭数据库连接
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("数据库连接失败")
		panic(err)
	}
	sqlDB.Close()
}

// 连接redis
func NewRedis(rds *Redis) (redis.Conn, error) {
	// 连接redis
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", rds.Host, rds.Port))
	if err != nil {
		logger.Errorf("redis连接失败")
		panic(err)
	}
	// 密码
	_, err = conn.Do("AUTH", rds.Password)
	if err != nil {
		logger.Errorf("redis密码错误")
		panic(err)
	}
	return conn, err
}

// 构建reddis连接池
func NewRedisPool(rds *Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,  //最大空闲连接数
		MaxActive:   0,   //最大连接数,0表示没有限制
		IdleTimeout: 600, //最大空闲时间
		Dial: func() (redis.Conn, error) {
			return NewRedis(rds)
		},
	}
}
