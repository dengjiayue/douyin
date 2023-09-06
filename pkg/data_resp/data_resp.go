package data_resp

//与提供返回数据不相同,暂时不要使用(data不一样)

// 封装统一的返回数据结构
type Resp struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg"`
	Data       interface{} `json:"data"`
}

// 返回成功
func Ok(data interface{}) *Resp {
	return &Resp{
		StatusCode: 200,
		StatusMsg:  "success",
		Data:       data,
	}
}

// 返回失败
func Fail(code int, msg string) *Resp {
	return &Resp{
		StatusCode: code,
		StatusMsg:  msg,
		Data:       nil,
	}
}
