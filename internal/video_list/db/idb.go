package db

import (
	"douyin/api/pb/user"
	"douyin/api/pb/video_list"
	"douyin/internal/user/service/client"
	"douyin/internal/video_list/model"
	"douyin/pkg/db"
	"douyin/pkg/logger"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	jsoner "github.com/json-iterator/go"
	"gorm.io/gorm"
)

// 定义接口
type IVideoListDB interface {
	//获取视频feed流列表,最多三十条,按时间先后倒序排列
	Feed(req *video_list.RpcFeedRequest) (resp *video_list.DouyinFeedResponse, err error)
	//关闭redis连接
	Close()
	//数据同步函数
	Sync(minTime int64)
}

// 定义结构体  TODO: 传入 gorm.DB 和 redis.DB
type VideoListDB struct {
	db   *gorm.DB
	pool *redis.Pool
	//继承user_client接口
	client.IUserClient
}

// 实现new方法
func NewVideoListDB(msq *db.Mysql, rds *db.Redis) (videoListDB *VideoListDB) {
	videoListDB = &VideoListDB{db.NewDB(msq), db.NewRedisPool(rds), client.NewUserClient()}
	//数据同步,当前时间戳
	go videoListDB.Sync(time.Now().Unix())
	return
}

// 通过时间同步
// 数据同步函数,小于某个时间戳的数据
func (videoListDB *VideoListDB) Sync(minTime int64) {
	//获取redis连接
	conn := videoListDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//将mysql的值缓存到redis中
	//按时间倒序输出
	var data []model.VideoData
	var err error
	videoListDB.db.Select("id, play_url, cover_url, favorite_count, comment_count, is_favorite, title,create_time,user_id").
		Where("create_time < ?", minTime).
		Order("create_time desc").
		Limit(100).
		Find(&data)
	lenth := len(data) * 2
	//校验数据长度
	if lenth == 0 {
		logger.Debugf("无需同步数据")
		return
	}
	//拼接数据interface切片
	videoData := make([]interface{}, lenth)

	video_time := make([]interface{}, lenth+1)
	video_time[0] = "video_time"
	for i := 0; i < lenth; i += 2 {
		video_time[i+1] = (data[i/2].CreateTime)
		video_time[i+2] = data[i/2].Id
		videoData[i] = fmt.Sprintf("video:%d", data[i/2].Id)
		videoData[i+1], err = jsoner.Marshal(data[i/2])
		if err != nil {
			logger.Errorf("json序列化失败: %v", err)
		}
	}
	//将数据同步到redis中
	_, err = conn.Do("MSET", videoData...)
	if err != nil {
		logger.Errorf("redis同步数据失败: %v", err)
	}
	_, err = conn.Do("ZADD", video_time...)
	if err != nil {
		logger.Errorf("redis同步数据失败: %v", err)
	}

	//数据同步到redis完成
	logger.Debugf("数据同步到redis完成")
}

// 通过datas同步数据
func SyncByDatas(videoListDB *VideoListDB, datas []model.VideoData) {
	//获取redis连接
	conn := videoListDB.pool.Get()
	//放回连接池
	defer conn.Close()

	//构建interface切片
	lenth := len(datas) * 2
	//校验数据长度
	if lenth == 0 {
		logger.Debugf("无需同步数据")
		return
	}
	//拼接数据interface切片
	var err error
	videoData := make([]interface{}, lenth)
	for i := 0; i < lenth; i += 2 {
		videoData[i] = fmt.Sprintf("video:%d", datas[i/2].Id)
		videoData[i+1], err = jsoner.Marshal(datas[i/2])
		if err != nil {
			logger.Errorf("json序列化失败: %v", err)
			return
		}
	}
	//将数据同步到redis中
	_, err = conn.Do("MSET", videoData...)
	if err != nil {
		logger.Errorf("redis同步数据失败: %v", err)
	}
}

// 关闭连接
func (videoListDB *VideoListDB) Close() {
	videoListDB.pool.Close()
	db.CloseDB(videoListDB.db)
}

// 批量获取视频数据
func VideosInfo(videoListDB *VideoListDB, ids []int64) []model.VideoData {
	//获取redis连接
	conn := videoListDB.pool.Get()
	//放回连接池
	defer conn.Close()
	lenth := len(ids)
	//校验数据长度
	if lenth == 0 {
		logger.Errorf("数据为空,无需查询")
		return nil
	}
	videos := make([]model.VideoData, lenth)
	//拼接数据interface切片
	data := make([]interface{}, lenth)

	for i := 0; i < lenth; i++ {
		data[i] = fmt.Sprintf("video:%d", ids[i])
	}
	//记录需要从mysql中获取的数量
	n := 0
	//从redis中获取数据
	videoJsons, err := redis.ByteSlices(conn.Do("MGET", data...))
	if err != nil {
		logger.Errorf("redis获取数据失败: %v", err)
		n = lenth
	} else {
		//遍历数据
		for i := 0; i < lenth; i++ {
			{
				//判断数据是否为空
				if videoJsons[i] == nil || len(videoJsons[i]) == 0 {
					videos[i].CreateTime = -1
					//将需要再次从mysql中获取的数据id记录下来
					ids[n] = ids[i]
					n++
				} else {
					err := jsoner.Unmarshal(videoJsons[i], &videos[i])
					if err != nil {
						logger.Errorf("json反序列化失败: %v", err)
						continue
					}
				}

			}
		}
	}
	sqlData := make([]model.VideoData, n)
	//从mysql中获取redis中没有的数据
	err = videoListDB.db.Where("id in ?", ids[:n]).Find(&sqlData).Error
	if err != nil {
		logger.Errorf("mysql获取数据失败: %v", err)
	}
	//将mysql中获取的数据同步到redis中
	go SyncByDatas(videoListDB, sqlData)
	//将MySQL查询到的数据合并到users中(为nil的就合并数据进去)
	for i, j := 0, 0; i < lenth && j < n; i++ {
		if videos[i].CreateTime == -1 {
			videos[i] = sqlData[j]
			j++
		}
	}
	return videos
}

// 获取视频feed流列表,最多三十条,按时间先后倒序排列
func (videoListDB *VideoListDB) Feed(req *video_list.RpcFeedRequest) (resp *video_list.DouyinFeedResponse, err error) {
	//获取redis连接
	conn := videoListDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//收到请求
	logger.Debugf("收到请求:%#v", req.LatestTime)
	//从redis中获取数据
	//获取视频id,从有序集合中获取,根据req.LatestTime的时间取小于这个时间的数据30条
	//如果没有传入时间,则默认取最新的30条数据
	var ids []int64

	//从redis中获取视频id
	if req.LatestTime == 0 {
		ids, err = redis.Int64s(conn.Do("ZREVRANGE", "video_time", 0, 29))
		req.LatestTime = time.Now().Unix()
	} else {
		ids, err = redis.Int64s(conn.Do("ZREVRANGEBYSCORE", "video_time", req.LatestTime, 0, "LIMIT", 0, 29))
	}
	var data []model.VideoData
	if err != nil || len(ids) == 0 {
		logger.Errorf("redis获取视频id失败:%v", err)

		//从mysql中获取数据(获取LatestTime>?的30条数据,并且按时间倒序排列)\

		err = videoListDB.db.Select("id, play_url, cover_url, favorite_count, comment_count, is_favorite, title,create_time,user_id").
			Where("create_time > ?", req.LatestTime).
			Order("create_time desc").
			Limit(30).Find(&data).Error
		if err != nil {
			logger.Errorf("mysql获取数据失败:%v", err)
			return &video_list.DouyinFeedResponse{StatusCode: 500, StatusMsg: "服务器出错"}, err
		}
		//将数据同步到redis中
		go SyncByDatas(videoListDB, data)
	} else {
		//获取视频数据
		data = VideosInfo(videoListDB, ids)
	}

	//如果没有数据,则返回空
	if data[len(data)-1].CreateTime == req.LatestTime {
		//如果最后一条数据的时间和请求的时间相同,则返回空
		return &video_list.DouyinFeedResponse{StatusCode: 0, StatusMsg: "没有更多视频了", VideoList: []*video_list.Video{}, NextTime: 0}, nil
	}

	//对鉴权用户: 获取点赞的视频id map
	favoriteVideoIds := make(map[int64]bool)
	if req.UserId != 0 {
		favoriteVideoIds, err = videoListDB.GetFavoriteVideoIds(req.UserId, ids)
		if err != nil {
			logger.Errorf("获取点赞的视频id失败:%v", err)
			return &video_list.DouyinFeedResponse{StatusCode: 500, StatusMsg: "服务器出错"}, err
		}
	}

	if len(ids) < 30 {
		//去数据库同步数据(以最后一条数据为开始时间同步100条数据)
		go videoListDB.Sync(data[len(data)-1].CreateTime)
	}
	//获取视频作者id
	authorIds := []int64{}
	//去重
	for i, authors := 0, make(map[int64]bool); i < len(data); i++ {
		if !authors[data[i].User_id] {
			authorIds = append(authorIds, data[i].User_id)
			authors[data[i].User_id] = true
		}
	}

	//检查视频信息
	// logger.Debugf("检查视频信息:%#v", data)
	//查询视频作者信息  :调user的rpc接口,然后查询
	usersData, err := VideoListDB.UsersInfo(*videoListDB, &user.DouyinUsersRequest{UserId: req.UserId, UserIds: authorIds})
	if err != nil {
		logger.Errorf("rpc调用失败:%v", err)
		return &video_list.DouyinFeedResponse{StatusCode: 500, StatusMsg: "服务器出错"}, err
	}
	//建立哈希表(用户id和用户信息的映射)
	users := make(map[int64]*user.User)
	for i := 0; i < len(usersData.Users); i++ {
		users[usersData.Users[i].Id] = usersData.Users[i]
	}

	//将用户信息和视频信息合并
	videos := make([]*video_list.Video, len(data))
	for i := 0; i < len(data); i++ {
		videos[i] = (data)[i].Video
		videos[i].VideoAuthor = users[(data)[i].User_id]
		videos[i].IsFavorite = favoriteVideoIds[(data)[i].Id]
	}
	//返回数据
	resp = &video_list.DouyinFeedResponse{StatusCode: 0, StatusMsg: "success", VideoList: videos, NextTime: data[len(data)-1].CreateTime}
	return
}

// 通过video_id获取自己点赞的视频id map<video_id>bool
func (videoListDB *VideoListDB) GetFavoriteVideoIds(userId int64, videoids []int64) (videoIds map[int64]bool, err error) {
	//通过user_id从数据库中获取视频id列表
	var video_ids []int64
	if err := videoListDB.db.Model(&model.Favorite{}).Where("user_id = ? AND is_favorite = ? AND video_id IN (?)", userId, 1, videoids).Pluck("video_id", &video_ids).Error; err != nil {
		return nil, err
	}
	videoIds = make(map[int64]bool, len(video_ids))
	for _, video_id := range video_ids {
		videoIds[video_id] = true
	}
	return
}
