package db

import (
	"douyin/api/pb/user"
	"douyin/api/pb/video"
	"douyin/api/pb/video_list"
	"douyin/internal/user/service/client"
	video_mod "douyin/internal/video/model"
	"douyin/internal/video_list/model"
	"fmt"

	jsoner "github.com/json-iterator/go"

	"douyin/pkg/db"
	"douyin/pkg/logger"
	"douyin/pkg/my_cos"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gorm.io/gorm"
)

// 定义接口
type IVideoDB interface {
	//上传视频publish/action
	PublishAction(req video.Video_DouyinPublishActionServer) (err error)
	//获取视频发布列表publish/list
	PublishList(req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error)
	//关闭redis连接
	Close()
	//点赞视频favorite/action
	FavoriteAction(req *video.RpcFavoriteActionRequest) (resp *video.DouyinFavoriteActionResponse, err error)
	//获取视频点赞列表favorite/list
	FavoriteList(req *video.DouyinFavoriteListRequest) (resp *video.DouyinFavoriteListResponse, err error)
	//评论操作comment/action
	CommentAction(req *video.RpcCommentActionRequest) (resp *video.DouyinCommentActionResponse, err error)
	//获取评论列表comment/list
	CommentList(req *video.RpcCommentListRequest) (resp *video.DouyinCommentListResponse, err error)
}

// 定义结构体
type VideoDB struct {
	//数据库连接
	db   *gorm.DB
	pool *redis.Pool
	//腾讯云cos连接
	client *cos.Client
	//继承user_client接口
	client.IUserClient
}

var _ IVideoDB = (*VideoDB)(nil)

// 实现new方法
func NewVideoDB(msq *db.Mysql, rds *db.Redis) (videoDB *VideoDB) {
	videoDB = &VideoDB{db.NewDB(msq), db.NewRedisPool(rds), my_cos.NewCosClient(), client.NewUserClient()}
	return
}

// 关闭redis连接
func (videoDB *VideoDB) Close() {
	videoDB.pool.Close()
	db.CloseDB(videoDB.db)
}

// 通过datas同步数据
func SyncByDatas(videoDB *VideoDB, datas []model.VideoData) {
	//获取redis连接
	conn := videoDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//构建interface切片
	lenth := len(datas) * 2
	//校验数据长度
	if lenth == 0 {
		logger.Debugf("数据为空,无需同步")
		return
	}
	//构建interface切片
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
		return
	}
}

// 批量获取视频数据
func VideosInfo(videoDB *VideoDB, ids *[]int64) *[]model.VideoData {
	//获取redis连接
	conn := videoDB.pool.Get()
	//放回连接池
	defer conn.Close()
	lenth := len(*ids)
	//校验数据长度
	if lenth == 0 {
		logger.Debugf("数据为空,无需查询")
		return nil
	}
	videos := make([]model.VideoData, lenth)
	//拼接数据interface切片
	data := make([]interface{}, lenth)

	for i := 0; i < lenth; i++ {
		data[i] = (*ids)[i]
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
					//将需要再次从mysql中获取的数据id记录下来
					(*ids)[n] = (*ids)[i]
					n++
				}

				err := jsoner.Unmarshal(videoJsons[i], &videos[i])
				if err != nil {
					continue
				}
			}
		}
	}
	sqlData := make([]model.VideoData, 0, n)
	//从mysql中获取redis中没有的数据
	err = videoDB.db.Where("id in ?", (*ids)[:n]).Find(&sqlData).Error
	if err != nil {
		logger.Errorf("mysql获取数据失败: %v", err)
	}
	//将mysql中获取的数据同步到redis中
	go SyncByDatas(videoDB, sqlData)
	//将MySQL查询到的数据合并到users中(为nil的就合并数据进去)
	for i, j := 0, 0; i < lenth && j < len(sqlData); i++ {
		if videos[i].User_id == 0 {
			videos[i] = sqlData[j]
			j++
		}
	}
	return &videos
}

// 同步redis中的赞数到mysql中(redis为赞的增量,总数=redis+mysql)
//考虑到数据及时性问题,暂时不采用定时任务,而是在点赞和取消点赞的时候同步
// func SyncFavorite(videoDB *VideoDB) {
// 	//从redis中获取视频id和赞数
// 	vids, err := redis.Int64Map(videoDB.conn.Do("HGETALL", "video_favorite"))
// 	if err != nil {
// 		return
// 	}

// }

// 上传视频publish/action
// 1.上传视频到腾讯云cos(视频上传cos-> 获取视频封面+上传封面-> 获取视频播放地址+封面地址)-> 2.创建视频结构体数据-> 3.将视频信息存入数据库-> 4.将视频信息存入redis
func (videoDB *VideoDB) PublishAction(req video.Video_DouyinPublishActionServer) (err error) {
	//获取redis连接
	conn := videoDB.pool.Get()
	//放回连接池
	defer conn.Close()
	//先接收用户与title数据
	header, err := req.Recv()
	if err != nil {
		logger.Errorf("接收用户与title数据失败: %v", err)
		return err
	}
	//获取当前时间戳(毫秒)
	creat_time := time.Now().UnixMilli()
	//上传视频到腾讯云cos
	Vurl, Purl := my_cos.UploadFile(videoDB.client, req, &header.UserId)
	if len(Vurl) == 0 {
		logger.Errorf("上传视频到腾讯云cos失败: %v", err)
		return err
	}

	//创建视频数据
	videoData := &model.VideoData{User_id: header.UserId, CreateTime: creat_time, Video: &video_list.Video{Title: header.Title, PlayUrl: Vurl, CoverUrl: Purl}}
	//将视频信息存入数据库
	if err := videoDB.db.Create(videoData).Error; err != nil {
		logger.Errorf("将视频信息存入数据库失败: %v", err)
		return err
	} else {
		//将视频信息存入redis有序集合中,时间戳作为score,视频id作为value

		_, err = conn.Do("ZADD", "video_time", creat_time, videoData.Id)
		if err != nil {
			logger.Errorf("将视频信息存入redis有序集合中失败: %v", err)
			return err
		}
		//调用user的rpc服务,修改用户的视频数量
		if _, err := videoDB.IUserClient.UserChange(&user.DouyinUserChangeRequest{UserId: header.UserId, ToUserId: 0, Type: 4, ActionType: 1}); err != nil {
			logger.Errorf("调用user的rpc服务,修改用户的视频数量失败: %v", err)
			return err
		}

		logger.Debugf("成功")
		return nil
	}
}

// 获取视频发布列表publish/list
// 网关鉴权->video服务:根据token获取用户id,根据用户id获取从mysql中获取视频列表->user服务:调用user服务获取用户信息->video服务:将用户信息和视频列表组合成响应数据返回给网关
func (videoDB *VideoDB) PublishList(req *video.DouyinPublishListRequest) (resp *video.DouyinPublishListResponse, err error) {

	//从mysql中获取视频列表
	var videodatas []*model.VideoData
	if err := videoDB.db.Where("user_id = ?", req.UserId).Find(&videodatas).Error; err != nil {
		logger.Errorf("从mysql中获取视频列表失败: %v", err)
		return &video.DouyinPublishListResponse{StatusMsg: "从mysql中获取视频列表失败", StatusCode: 500}, err
	}
	//获取视频列表ids
	var video_ids []int64
	for i := 0; i < len(videodatas); i++ {
		video_ids = append(video_ids, videodatas[i].Id)
	}
	//获取点赞视频map
	favoriteVideoIds, err := videoDB.GetFavoriteVideoIds(req.UserId, video_ids)
	if err != nil {
		logger.Errorf("获取点赞视频map失败: %v", err)
		return &video.DouyinPublishListResponse{StatusMsg: "获取点赞视频map失败", StatusCode: 500}, err
	}

	//调用user服务获取用户信息
	user_info, err := videoDB.UserInfo(&user.DouyinUserRequest{UserId: req.UserId})
	if err != nil {
		logger.Errorf("调用user服务获取用户信息失败: %v", err)
		return &video.DouyinPublishListResponse{StatusMsg: "调用user服务获取用户信息失败", StatusCode: 500}, err
	}
	//视频列表
	video_list := make([]*video_list.Video, len(videodatas))
	//将用户信息和视频列表组合成响应数据返回给网关
	for i := 0; i < len(videodatas); i++ {
		video_list[i] = videodatas[i].Video
		video_list[i].VideoAuthor = user_info.User
		video_list[i].IsFavorite = favoriteVideoIds[videodatas[i].Id]
	}
	//将视频列表和用户信息组合成响应数据返回给网关
	return &video.DouyinPublishListResponse{StatusMsg: "获取视频列表成功", StatusCode: 0, VideoList: video_list}, nil
}

// 点赞视频favorite/action
// 一赞一存,每一个赞操作,先储存sql,再更新redis相关内容(数据操作包括:点赞表,视频表,用户表的数据更新,以及对视频与用户的redis数据进行更新)
func (videoDB *VideoDB) FavoriteAction(req *video.RpcFavoriteActionRequest) (resp *video.DouyinFavoriteActionResponse, err error) {
	//获取redis连接
	conn := videoDB.pool.Get()
	//放回连接池
	defer conn.Close()

	//检查数据库favorite表中是否的点赞数据是否与请求的点赞数据一致,如果一致就不进行操作,如果不一致就进行操作
	favorite := video_mod.Favorite{VideoId: req.VideoId, UserId: req.UserId, IsFavorite: req.ActionType}

	//查询/创建点赞记录
	result := videoDB.db.Where("user_id = ? AND video_id = ?", req.UserId, req.VideoId).FirstOrCreate(&favorite)

	//非法操作(未关注时,取消关注)
	if result.RowsAffected == 1 && req.ActionType == 2 {
		logger.Errorf("非法操作:未点赞取消点赞")
		return &video.DouyinFavoriteActionResponse{StatusMsg: "非法操作:未点赞取消点赞", StatusCode: 500}, fmt.Errorf("非法操作")
	}

	//是否需要更新(存在值并且与传入值不同)
	if favorite.IsFavorite != req.ActionType {
		//更新点赞记录
		result = videoDB.db.Model(&video_mod.Favorite{}).Where("user_id = ? AND video_id = ?", req.UserId, req.VideoId).Update("is_favorite", req.ActionType)
		result.RowsAffected = 1
	}

	//判断是否需要更改数据(改变行数=0,返回成功,改变行数=1,继续修改数据)
	if result.RowsAffected == 0 {
		logger.Errorf("已经点赞过了")
		return &video.DouyinFavoriteActionResponse{StatusMsg: "已经点赞过了", StatusCode: 0}, nil
	}

	//--------------继续修改数据----------------
	// 开启事务
	tx := videoDB.db.Begin()
	defer tx.Commit()
	//更改video_list表中的点赞数
	if req.ActionType == 1 {
		if err := videoDB.db.Model(&model.VideoData{}).Where("id = ?", req.VideoId).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			logger.Errorf("更改video_list表中的点赞数失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}
		//查询视频的作者id
		var videoData *model.VideoData
		if err := videoDB.db.Where("id = ?", req.VideoId).First(&videoData).Error; err != nil {
			logger.Errorf("查询视频的作者id失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}
		//更新redis中的视频数据
		//json序列化
		video_json, err := jsoner.Marshal(videoData.Video)
		if err != nil {
			logger.Errorf("json序列化失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}
		//更新redis中的视频数据
		_, err = conn.Do("SET", fmt.Sprintf("video:%d", req.VideoId), video_json)
		if err != nil {
			logger.Errorf("更新redis中的视频数据失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}

		//调user的rpc服务,修改用户的点赞数量
		if _, err := videoDB.IUserClient.UserChange(&user.DouyinUserChangeRequest{UserId: req.UserId, ToUserId: videoData.User_id, Type: 0, ActionType: req.ActionType}); err != nil {
			logger.Errorf("调user的rpc服务,修改用户的点赞数量失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}

	} else if req.ActionType == 2 {
		if err := videoDB.db.Model(&model.VideoData{}).Where("id = ?", req.VideoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			logger.Errorf("更改video_list表中的点赞数失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "取消点赞失败", StatusCode: 500}, err
		}
		//查询视频的作者id
		var videoData *model.VideoData
		if err := videoDB.db.Where("id = ?", req.VideoId).First(&videoData).Error; err != nil {
			logger.Errorf("查询视频的作者id失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞失败", StatusCode: 500}, err
		}
		//更新redis中的视频数据
		//json序列化
		video_json, err := jsoner.Marshal(videoData.Video)
		if err != nil {
			logger.Errorf("json序列化失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "取消点赞失败", StatusCode: 500}, err
		}
		//更新redis中的视频数据
		_, err = conn.Do("SET", fmt.Sprintf("video:%d", req.VideoId), video_json)
		if err != nil {
			logger.Errorf("更新redis中的视频数据失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "取消点赞失败", StatusCode: 500}, err
		}

		//调user的rpc服务,修改用户的点赞数量
		if _, err := videoDB.IUserClient.UserChange(&user.DouyinUserChangeRequest{UserId: req.UserId, ToUserId: videoData.User_id, Type: 0, ActionType: req.ActionType}); err != nil {
			logger.Errorf("调user的rpc服务,修改用户的点赞数量失败: %v", err)
			tx.Rollback()
			return &video.DouyinFavoriteActionResponse{StatusMsg: "取消点赞失败", StatusCode: 500}, err
		}

	}
	//提交事务
	tx.Commit()
	return &video.DouyinFavoriteActionResponse{StatusMsg: "点赞成功", StatusCode: 0}, nil
}

// 获取视频点赞列表favorite/list
// 数据流动: 网关(鉴权)-video服务(redis查询视频id->批量查询视频(查不到的id记录下来->sql查询redis查不到的数据)获取视频用户id)->用户服务(批量查询用户信息)->video服务(组装数据返回)
func (videoDB *VideoDB) FavoriteList(req *video.DouyinFavoriteListRequest) (resp *video.DouyinFavoriteListResponse, err error) {

	//通过user_id从数据库中获取视频id列表
	var video_ids []int64
	if err := videoDB.db.Model(&video_mod.Favorite{}).Where("user_id = ? AND is_favorite = ?", req.UserId, 1).Pluck("video_id", &video_ids).Error; err != nil {
		logger.Errorf("从mysql中获取视频id列表失败: %v", err)
		return &video.DouyinFavoriteListResponse{StatusMsg: "从mysql中获取视频id列表失败", StatusCode: "500"}, err
	}
	// 检查数据长度
	if len(video_ids) == 0 {
		logger.Debugf("数据为空,无需查询")
		return &video.DouyinFavoriteListResponse{StatusMsg: "无点赞视频,数据为空", StatusCode: "200"}, nil
	}
	//批量查询视频信息
	videosData := VideosInfo(videoDB, &video_ids)
	//检查数据
	if videosData == nil || len(*videosData) == 0 {
		logger.Debugf("数据为空,无需查询")
		return &video.DouyinFavoriteListResponse{StatusMsg: "获取视频数据为空", StatusCode: "500"}, nil
	}

	//获取视频列表的用户ids
	var user_ids []int64
	for i := 0; i < len(*videosData); i++ {
		user_ids = append(user_ids, (*videosData)[i].User_id)

	}

	//批量查询用户信息
	user_list, err := videoDB.UsersInfo(&user.DouyinUsersRequest{UserIds: user_ids})
	if err != nil {
		logger.Errorf("批量查询用户信息失败: %v", err)
		return &video.DouyinFavoriteListResponse{StatusMsg: "批量查询用户信息失败", StatusCode: "500"}, err
	}
	videos := make([]*video_list.Video, len(*videosData))
	//将用户信息和视频列表组合成响应数据返回给网关
	for i, j := 0, 0; i < len(*videosData); i++ {
		//如果视频作者id和用户id相同,就将用户信息赋值给视频作者信息
		if (*videosData)[i].User_id == user_list.Users[j].Id && j < len(user_list.Users) && user_list.Users[j].Id != 0 {
			(*videosData)[i].VideoAuthor = user_list.Users[j]
			j++
		}
		videos[i] = (*videosData)[i].Video
		//标记为已点赞(条件:不为空)
		if videos[i] != nil {
			videos[i].IsFavorite = true
		}
	}
	//将视频列表和用户信息组合成响应数据返回给网关
	return &video.DouyinFavoriteListResponse{StatusMsg: "获取视频列表成功", StatusCode: "0", VideoList: videos}, nil
}

// 评论操作comment/action(1评论,2删除评论)
//   - 数据流动:网关(用户鉴权)->视频服务: 评论/取消评论(1. 将评论信息存入/删除数据库对应数据 ,2. 更改视频的评论数(评论数+1/-1), 3. 清除redis的视频缓存(清除旧数据缓存))->网关->响应
func (videoDB *VideoDB) CommentAction(req *video.RpcCommentActionRequest) (resp *video.DouyinCommentActionResponse, err error) {
	//获取redis连接
	conn := videoDB.pool.Get()
	//放回连接池
	defer conn.Close()

	//1评论,2删除评论
	data := &video_mod.CommentData{UserId: req.UserId, VideoId: req.VideoId, Content: req.CommentText, CreateTime: time.Now().Unix()}
	// 开启事务
	tx := videoDB.db.Begin()
	if req.ActionType == 1 {
		//向数据库中存入数据
		if err := videoDB.db.Create(data).Error; err != nil {

			tx.Rollback()
			logger.Errorf("评论失败,err: %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "评论失败", StatusCode: 500}, err
		}
		//更改video_list表中的视频的评论数
		if err := videoDB.db.Model(&model.VideoData{}).Where("id = ?", req.VideoId).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			logger.Errorf("评论失败,err: %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "评论失败", StatusCode: 500}, err
		}
		//清除redis中的视频数据
		if _, err := conn.Do("DEL", fmt.Sprintf("video:%d", req.VideoId)); err != nil {
			tx.Rollback()
			logger.Errorf("评论失败,err: %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "评论失败", StatusCode: 500}, err
		}
		//提交事务
		tx.Commit()
		logger.Debugf("评论成功,data: %v", data)
		return &video.DouyinCommentActionResponse{StatusMsg: "评论成功", StatusCode: 0, Comment: data.ToComment()}, nil
	} else if req.ActionType == 2 {
		//检查参数
		if req.CommentId == 0 {
			tx.Rollback()
			logger.Errorf("删除评论失败:err= %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "参数错误", StatusCode: 500}, fmt.Errorf("删除评论失败")
		}

		//删除评论(commentid与 userid(判断是否是自己的评论)同时匹配才能删除)
		if n := videoDB.db.Where("user_id = ?", req.UserId).Delete(&video_mod.CommentData{Id: req.CommentId}).RowsAffected; n == 0 {
			tx.Rollback()
			logger.Errorf("删除评论失败:err= %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "删除评论失败", StatusCode: 500}, err
		}
		//更改video_list表中的视频的评论数
		if err := videoDB.db.Model(&model.VideoData{}).Where("id = ?", req.VideoId).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			logger.Errorf("删除评论失败:err= %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "删除评论失败", StatusCode: 500}, err
		}
		//清除redis中的视频数据
		if _, err := conn.Do("DEL", fmt.Sprintf("video:%d", req.VideoId)); err != nil {
			tx.Rollback()
			logger.Errorf("删除评论失败:err= %v", err)
			return &video.DouyinCommentActionResponse{StatusMsg: "删除评论失败", StatusCode: 500}, err
		}
		//提交事务
		tx.Commit()
		logger.Debugf("删除评论成功,data: %v", data)
		return &video.DouyinCommentActionResponse{StatusMsg: "删除评论成功", StatusCode: 0}, nil
	}
	return &video.DouyinCommentActionResponse{StatusMsg: "参数错误", StatusCode: 500}, fmt.Errorf("删除评论失败")
}

// 评论列表comment/list
//
//	数据流动:网关(用户鉴权)->视频服务: 获取评论列表(1. 从数据库获取评论列表(按时间倒序),取出评论列表的使用user的ids)->用户服务: 批量查询用户数据(批量查询用户数据)->视频服务(数据整合,将user数据与comment数据合并返回)->网关->响应
func (videoDB *VideoDB) CommentList(req *video.RpcCommentListRequest) (resp *video.DouyinCommentListResponse, err error) {
	//通过video_id从数据库中获取评论列表
	var comments []video_mod.CommentData
	if err := videoDB.db.Where("video_id = ?", req.VideoId).Order("create_time desc").Find(&comments).Error; err != nil {
		logger.Errorf("从mysql中获取评论列表失败: %v", err)
		return &video.DouyinCommentListResponse{StatusMsg: "从mysql中获取评论列表失败", StatusCode: 500}, err
	}
	//获取评论列表的用户ids
	var user_ids []int64
	for i := 0; i < len(comments); i++ {
		user_ids = append(user_ids, comments[i].UserId)
	}
	//批量查询用户信息
	user_list, err := videoDB.UsersInfo(&user.DouyinUsersRequest{UserId: req.UserId, UserIds: user_ids})
	if err != nil {
		logger.Errorf("批量查询用户信息失败: %v", err)
		return &video.DouyinCommentListResponse{StatusMsg: "批量查询用户信息失败", StatusCode: 500}, err
	}
	//将用户信息和评论列表组合成响应数据返回给网关
	comment_list := make([]*video.Comment, len(comments))
	//如果user的id与comment的user_id相同,就将user信息赋值给comment的user信息
	for i, j := 0, 0; i < len(comments); i++ {
		comment_list[i] = comments[i].ToComment()
		if comments[i].UserId == user_list.Users[j].Id {
			comment_list[i].User = user_list.Users[j]
			j++
		}
	}
	logger.Debugf("获取评论列表成功,comment_list: %v", comment_list)
	return &video.DouyinCommentListResponse{StatusMsg: "获取评论列表成功", StatusCode: 0, CommentList: comment_list}, nil
}

// 通过video_id获取自己点赞的视频id map<video_id>bool
func (videoDB *VideoDB) GetFavoriteVideoIds(userId int64, videoids []int64) (videoIds map[int64]bool, err error) {
	//通过user_id从数据库中获取视频id列表
	var video_ids []int64
	if err := videoDB.db.Model(&video_mod.Favorite{}).Where("user_id = ? AND is_favorite = ? AND video_id IN (?)", userId, 1, videoids).Pluck("video_id", &video_ids).Error; err != nil {
		logger.Errorf("从mysql中获取视频id列表失败: %v", err)
		return nil, err
	}
	videoIds = make(map[int64]bool, len(video_ids))
	for _, video_id := range video_ids {
		videoIds[video_id] = true
	}
	return
}
