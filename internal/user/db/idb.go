package db

import (
	"douyin/api/pb/social"
	"douyin/api/pb/user"
	"douyin/internal/social/service/client"
	"douyin/internal/user/model"
	"douyin/pkg/db"
	"douyin/pkg/logger"
	"fmt"
	"math/rand"

	"github.com/gomodule/redigo/redis"
	jsoner "github.com/json-iterator/go"
	"gorm.io/gorm"
)

var CHANGE_TYPE = []string{"favorite_count", "total_favorited", "follow_count", "follower_count", "work_count"} //字段解释：点赞数(你点别人)，获赞赞数(别人点你)，关注数(你关注别人)，粉丝数(别人关注你),作品数(你发布的作品数)
var ACTION_TYPE = []string{" + 1", " - 1"}                                                                      //字段解释： 1-增加操作,2-取消操作
var AVATAR = []string{"https://djy1-1306563712.cos.ap-shanghai.myqcloud.com/20230825005714.png"}                //默认头像
var BAKEGROUNDIMG = []string{"https://djy1-1306563712.cos.ap-shanghai.myqcloud.com/20230825005607.png"}         //默认背景图
var LENGTH = len(AVATAR)                                                                                        //头像长度

type IUserDB interface {
	// 注册用户
	RegisterUser(RegisterUserRespons *user.DouyinUserRegisterRequest) (respons *user.DouyinUserRegisterResponse, err error)
	// 登录
	LoginUser(LoginUserRespons *user.DouyinUserLoginRequest) (respons *user.DouyinUserLoginResponse, err error)
	// 查询用户信息
	UserInfo(UserInfoRespons *user.DouyinUserRequest) (respons *user.DouyinUserResponse, err error)
	//关闭ridis连接
	Close()
	//批量查询用户信息
	UsersInfo(UserInfosRespons *user.DouyinUsersRequest) (respons *user.DouyinUsersResponse, err error)
	//更新用户数据
	UserChange(UserChangeRespons *user.DouyinUserChangeRequest) (err error)
}

// UserDB 是用户数据库
type UserDB struct {
	db   *gorm.DB
	pool *redis.Pool
	//继承social_client接口
	client.ISocialClient
}

// 确保结构体实现接口的小技巧:
var _ IUserDB = (*UserDB)(nil)

// TODO: 传入 gorm.DB 和 redis.DB
func NewUserDB(sql *db.Mysql, rds *db.Redis) (userDB *UserDB) {
	userDB = &UserDB{db.NewDB(sql), db.NewRedisPool(rds), client.NewSocialClient()}

	//弃用:启动程序不再同步数据
	// //将mysql的值缓存到redis中
	// var userData []model.UserData
	// userDB.db.Find(&userData)
	// //批量同步数据到redis
	// SyncDataToRedis(userDB, &userData)
	return
}

// 批量同步数据到redis
func SyncDataToRedis(userDB *UserDB, userData *[]model.UserData) {
	//获取redis连接
	conn := userDB.pool.Get()
	defer conn.Close()
	//检查数据长度
	if len(*userData) == 0 {
		logger.Errorf("数据为空,无需同步")
		return
	}

	//拼接数据interface切片
	data := make([]interface{}, len(*userData)*2)

	for i := 0; i < len(*userData)*2; i += 2 {
		userJson, err := jsoner.Marshal((*userData)[i/2])
		if err != nil {
			logger.Errorf("json解析失败: %v", err)
			continue
		}
		data[i+1] = userJson
		data[i] = fmt.Sprintf("user:%d", (*userData)[i/2].Id)
	}
	//批量同步数据到redis
	conn.Do("MSET", data...)
}

// 批量获取用户信息
func Users(userDB *UserDB, ids *[]int64) []model.UserData {
	//获取redis连接
	conn := userDB.pool.Get()
	defer conn.Close()
	lenth := len(*ids)
	//校验数据长度
	if lenth == 0 {
		logger.Errorf("数据为空,无需查询")
		return nil
	}
	users := make([]model.UserData, lenth)
	//拼接数据interface切片
	data := make([]interface{}, lenth)
	for i := 0; i < lenth; i++ {
		data[i] = fmt.Sprintf("user:%d", (*ids)[i])
	}
	//记录需要从mysql中获取的数量
	n := 0
	//从redis中获取数据
	userJsons, err := redis.ByteSlices(conn.Do("MGET", data...))
	if err != nil {
		logger.Errorf("redis获取数据失败: %v", err)
		n = lenth
	} else {
		//遍历数据
		for i := 0; i < lenth; i++ {
			{
				//判断数据是否为空
				if userJsons[i] == nil || len(userJsons[i]) == 0 {

					//将需要再次从mysql中获取的数据id记录下来
					(*ids)[n] = (*ids)[i]
					n++
					continue
				}

				err := jsoner.Unmarshal(userJsons[i], &users[i])
				if err != nil {
					logger.Errorf("json反序列化失败: %v", err)
					//将需要再次从mysql中获取的数据id记录下来
					(*ids)[n] = (*ids)[i]
					n++
					continue
				}
			}
		}
	}
	sqlData := make([]model.UserData, n)
	//如果不为0，从mysql中获取数据
	if n == 0 {
		return users
	}
	//从mysql中获取redis中没有的数据
	err = userDB.db.Where("id in ?", (*ids)[:n]).Find(&sqlData).Error
	if err != nil {
		logger.Errorf("mysql获取数据失败: %v", err)
		return users
	}
	//将mysql中获取的数据同步到redis中
	go SyncDataToRedis(userDB, &sqlData)
	//将MySQL查询到的数据合并到users中(为nil的就合并数据进去)
	for i, j := 0, 0; i < len(users) && j < len(sqlData); i++ {
		if len(users[i].Password) == 0 {
			users[i] = sqlData[j]
			j++
		}
	}
	return users
}

// 关闭ridis连接
func (userDB *UserDB) Close() {
	userDB.pool.Close()
	db.CloseDB(userDB.db)
}

// 注册逻辑
func (userDB *UserDB) RegisterUser(req *user.DouyinUserRegisterRequest) (respons *user.DouyinUserRegisterResponse, err error) {

	//检查name
	logger.Debugf("用户注册data: %v", req)
	n := rand.Intn(LENGTH)
	userData := &model.UserData{Password: req.Password, User: &user.User{Name: req.Username, Avatar: AVATAR[n], BackgroundImage: BAKEGROUNDIMG[n], Signature: "云外一只喵"}}

	// 1. 查询用户是否存在
	err = userDB.db.Where("name = ?", req.Username).First(&userData).Error
	if err == nil {
		logger.Errorf("用户已经存在: %v", err)
		return &user.DouyinUserRegisterResponse{StatusCode: 500, StatusMsg: "用户已经存在"}, err
	} else {
		err = userDB.db.Create(&userData).Error
		if err != nil {
			logger.Errorf("创建用户失败: %v", err)
			return &user.DouyinUserRegisterResponse{StatusCode: 500, StatusMsg: "创建用户失败"}, err
		}
	}
	logger.Debugf("用户注册成功: %v", userData)
	respons = &user.DouyinUserRegisterResponse{StatusCode: 0, StatusMsg: "注册成功", UserId: userData.Id}
	//返回数据
	return respons, nil
}

// 登录逻辑
func (userDB *UserDB) LoginUser(req *user.DouyinUserLoginRequest) (respons *user.DouyinUserLoginResponse, err error) {
	// 1. 查询用户是否存在
	var userData model.UserData

	//从mysql中查询
	err = userDB.db.Where("name = ?", req.Username).First(&userData).Error
	if err != nil {
		logger.Errorf("用户不存在: %v", err)
		return &user.DouyinUserLoginResponse{StatusCode: 500, StatusMsg: "用户不存在"}, err
	}

	// 3. 用户存在，验证密码是否正确
	if userData.Password != req.Password {
		logger.Errorf("密码错误: %v", err)
		respons = &user.DouyinUserLoginResponse{StatusCode: 400, StatusMsg: "密码错误"}
		return respons, nil
	}
	// 4. 密码正确，返回用户id
	respons = &user.DouyinUserLoginResponse{StatusCode: 0, StatusMsg: "登录成功", UserId: userData.Id}
	//返回数据
	return respons, nil
}

// 查询用户信息:通过id,先查找redis，再查找mysql并缓存到redis
func (userDB *UserDB) UserInfo(req *user.DouyinUserRequest) (respons *user.DouyinUserResponse, err error) {
	//获取redis连接
	conn := userDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//检查参数
	logger.Debugf("用户信息查询data: %v", req)
	// 1. 查询用户是否存在
	var userData model.UserData

	userJSON, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("user:%d", req.UserId)))
	if err != nil {
		//redis中没有数据，从mysql中查询
		err = userDB.db.Where("id = ?", req.UserId).First(&userData).Error
		if err != nil {
			logger.Errorf("用户不存在: %v", err)
			return &user.DouyinUserResponse{StatusCode: 500, StatusMsg: "查询用户是否存在失败"}, err
		} else {
			//开协程处理
			go func() {
				//将mysql的值缓存到redis中
				userJson, err := jsoner.Marshal(userData)
				if err != nil {
					return
				}
				_, err = conn.Do("SET", fmt.Sprintf("user:%d", req.UserId), userJson)
				if err != nil {
					logger.Errorf("数据同步到redis失败: %v", err)
					return
				}
				//数据同步到redis完成
				logger.Debugf("数据同步到redis完成  data:%v", userData)
			}()
		}

	} else {
		//redis中有数据，解析json
		err = jsoner.Unmarshal(userJSON, &userData)
		if err != nil {
			logger.Errorf("用户不存在: %v", err)
			return &user.DouyinUserResponse{StatusCode: 500, StatusMsg: "查询用户是否存在失败"}, err
		}
	}

	// 3. 用户存在，返回用户信息
	respons = &user.DouyinUserResponse{StatusCode: 0, StatusMsg: "查询成功", User: userData.User}
	//返回数据
	return respons, nil
}

// 批量查询用户信息
func (userDB *UserDB) UsersInfo(req *user.DouyinUsersRequest) (respons *user.DouyinUsersResponse, err error) {
	//检查参数
	logger.Debugf("用户信息查询data: %v", req)
	// 批量获取用户信息
	userDatas := Users(userDB, &req.UserIds)
	users := make([]*user.User, len(userDatas))
	followMap := make(map[int64]bool)
	//如果用户id不为0，批量获取用户关注的用户
	if req.UserId != 0 {
		//批量获取用户关注的用户:id(map<int64>bool),在传入ids中查询用户关注的ids存map中
		followresp, err := userDB.FindFollows(&social.FindFollowsRequest{UserId: req.UserId, Ids: req.UserIds})
		if err != nil {
			logger.Errorf("查询失败: %v", err)
			return nil, err
		}
		followMap = followresp.FollowMap
	}
	for i := 0; i < len(userDatas); i++ {
		//判断用户是否存在
		if len(userDatas[i].Password) == 0 {
			logger.Errorf("用户不存在: %v", err)
			continue
		}
		users[i] = (userDatas)[i].User
		users[i].IsFollow = followMap[users[i].Id]
	}
	// 3. 用户存在，返回用户信息
	respons = &user.DouyinUsersResponse{Users: users}

	// 检查返回数据
	logger.Debugf("查询成功: %v", users)
	//返回数据
	return respons, nil
}

// 更新用户信息
func (userDB *UserDB) UserChange(req *user.DouyinUserChangeRequest) (err error) {
	//获取redis连接
	conn := userDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//检查参数
	logger.Debugf("用户信息更新data: %v", req)

	//更新MySQL中数据
	err = userDB.db.Model(&model.UserData{}).Where("id = ?", req.UserId).Update(CHANGE_TYPE[req.Type], gorm.Expr(CHANGE_TYPE[req.Type]+ACTION_TYPE[req.ActionType-1])).Error
	if err != nil {
		return err
	}

	//更新to_user用户的数据
	if req.ToUserId != 0 {
		err = userDB.db.Model(&model.UserData{}).Where("id = ?", req.ToUserId).Update(CHANGE_TYPE[req.Type+1], gorm.Expr(CHANGE_TYPE[req.Type+1]+ACTION_TYPE[req.ActionType-1])).Error
		if err != nil {
			return err
		}
	}

	// 清除redis中的缓存
	_, err = conn.Do("DEL", fmt.Sprintf("user:%d", req.UserId), fmt.Sprintf("user:%d", req.ToUserId))
	if err != nil {
		return err
	}
	return nil
}
