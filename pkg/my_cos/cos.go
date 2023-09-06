package my_cos

import (
	"bytes"
	"context"
	"douyin/api/pb/video"
	"douyin/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 定义枚举类型:cos地址(string)
const COS_URL = "https://douyin-1306563712.cos.ap-guangzhou.myqcloud.com/"

// 文件名逻辑:user_id:时间戳.后缀(图片:jpg,视频:mp4)
// 上传文件函数(接收grpc流式上传的视频数据bytes流类型,流式上传到cos,返回视频url和封面url)
func UploadFile(client *cos.Client, video video.Video_DouyinPublishActionServer, user_id *int64) (url string, picurl string) {
	// 获取当前时间戳(毫秒)
	t := time.Now().UnixMilli()
	//一个user一个文件夹
	url = fmt.Sprintf("user_id:%d/%d.mp4", *user_id, t)
	// 上传视频到cos
	//流式上传处理
	//创建分块对象
	header := &http.Header{}
	header.Add("Content-Type", "video/mp4")
	opt := &cos.MultiUploadOptions{
		OptIni: &cos.InitiateMultipartUploadOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				XOptionHeader: header,
			},
		},
	}

	//创建分块上传
	mu, _, err := client.Object.InitiateMultipartUpload(context.Background(), url, opt.OptIni)
	if err != nil {
		fmt.Println(err)
		return
	}

	partNumber := 1
	//分块上传
	var parts []cos.Object
	for {
		data, err := video.Recv()
		if err == io.EOF {
			break // 数据读取完毕
		}
		if err != nil {
			logger.Errorf("流式上传视频数据失败: %s\n", err)
			return "", ""
		}

		// partSize := len(data.Data)
		reader := bytes.NewReader(data.Data)
		resp, err := client.Object.UploadPart(context.Background(), url, mu.UploadID, partNumber, reader, &cos.ObjectUploadPartOptions{})
		if err != nil {
			logger.Errorf("流式上传视频数据失败: %s\n", err)
			return "", ""
		}
		parts = append(parts, cos.Object{PartNumber: partNumber, ETag: resp.Header.Get("ETag")})
		partNumber++
	}
	// 完成分块上传
	_, _, err = client.Object.CompleteMultipartUpload(
		context.Background(),
		url,
		mu.UploadID,
		&cos.CompleteMultipartUploadOptions{Parts: parts},
	)

	if err != nil {
		logger.Errorf("流式上传视频数据失败: %s\n", err)
		return "", ""
	}

	return COS_URL + url, fmt.Sprintf("%suser_id:%d/%d.jpg", COS_URL, *user_id, t)

	// // 拼接url
	// url = COS_URL + url
	// picurl = fmt.Sprintf("img/%d:%d.jpg", user_id, t)
	// //上传视频封面到cos
	// picData, err := GetVideoCover(*videoData)
	// if err != nil {
	// 	logger.Errorf("GetVideoCover failed, err:%v", err)
	// 	return url, ""
	// }
	// // 上传图片到cos
	// _, err = client.Object.Put(context.Background(), picurl, bytes.NewReader(picData), nil)
	// if err != nil {
	// 	logger.Errorf(" 图片上传失败: client.Object.Put failed, err:%v", err)
	// 	return url, ""
	// }
	// // 拼接url
	// picurl = COS_URL + picurl
	// return
}

//已弃用:使用cos的数据万象替代!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// // 截取视频封面:视频数据bytes的第一帧,返回图片的bytes数据
// func GetVideoCover(videoData []byte) (picData []byte, err error) {
// 	// 使用FFmpeg命令来提取第一帧图片数据，-vframes 1 表示只提取一帧
// 	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-vframes", "1", "-f", "image2", "-")
// 	cmd.Stdin = bytes.NewReader(videoData)
// 	var outBuffer, errBuffer bytes.Buffer
// 	cmd.Stdout = &outBuffer
// 	cmd.Stderr = &errBuffer

// 	err = cmd.Run()
// 	if err != nil {
// 		return nil, fmt.Errorf("FFmpeg error: %s, %s", err, errBuffer.String())
// 	}

// 	return outBuffer.Bytes(), nil
// }
