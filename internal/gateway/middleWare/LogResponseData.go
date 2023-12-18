package middleware

import (
	"bytes"
	"douyin/pkg/logger"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// 自定义 ResponseWriter 包装器
type ResponseWithLog struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseWithLog) Write(data []byte) (int, error) {
	r.body.Write(data)
	return r.ResponseWriter.Write(data)
}

// 中间件函数
func LogResponseDataMiddleware(c *gin.Context) {
	// 打印请求头部数据
	logger.Debugf("Request Headers:")
	for key, values := range c.Request.Header {
		logger.Debugf("%s: %s", key, values)
	}

	// 打印请求主体数据
	reqbody, _ := c.GetRawData()
	logger.Debugf("Request Body Size:%d\n", len(reqbody))

	size := 100
	if len(reqbody) < size {
		logger.Debugf("Request Body:%s\n", reqbody)
	} else {
		logger.Debugf("Request Body:%s\n", reqbody[:size])
	}

	// 重置请求主体数据
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqbody))

	// 创建自定义 ResponseWriter 包装器
	responseWithLog := &ResponseWithLog{
		ResponseWriter: c.Writer,
		body:           bytes.NewBufferString(""),
	}
	c.Writer = responseWithLog

	// 在处理请求之前执行的代码
	c.Next()

	// 在响应之后执行的代码
	status := c.Writer.Status()
	respbody := responseWithLog.body.String()

	// 打印响应数据
	logger.Debugf("Response Status: %d", status)
	logger.Debugf("Response Body: %s", respbody)
}
