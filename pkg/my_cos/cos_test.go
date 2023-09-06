package my_cos

import (
	"testing"
)

// -----------------停用此单元测试----------------
func TestNewCosClient(t *testing.T) {

	//初始化日志
	// logger.Init(nil)
	// client := NewCosClient()
	//图片数据:bytes字节流
	// // 将图片数据上传到 COS 存储桶中
	// r, err := client.Object.Put(context.Background(), "image/text.png", bytes.NewReader(imgdata), nil)
	// if err != nil {
	// 	fmt.Println("Upload image failed:", err)
	// } else {
	// 	fmt.Printf("Upload image succeeded r= %v", r)
	// }

	// //视频数据:
	// //上传视频到cos
	// v, p := UploadFile(client, &vdodata, 5)
	// fmt.Printf("purl=%s\n,vurl=%s\n", p, v)
	//获取第一帧图片数据
	// picdata, err := GetVideoCover(vdodata)
	// if err != nil {
	// 	fmt.Printf("err=%#v\n", err)
	// }
	// fmt.Printf("data=%#v\n", picdata)
}
