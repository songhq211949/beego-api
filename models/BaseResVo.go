package models

//BaseResVO 重构接口,基本返回对象
type BaseResVO struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//ResponseCode 返回状态码
type ResponseCode struct {
	Code    int
	Message string
}

var NOT_NETWORK = ResponseCode{Code: -1, Message: "系统繁忙，请稍后再试~"}
var SUCCESS = ResponseCode{Code: 0, Message: "success"}
var LOGIN_VERIFY_FALL = ResponseCode{Code: 1, Message: "登录失效~"}
var PARAM_VERIFY_FALL = ResponseCode{Code: 2, Message: "参数验证错误~"}
var AUTH_FAILED = ResponseCode{Code: 3, Message: "权限验证失败~"}
var DATA_NOT = ResponseCode{Code: 4, Message: "没有相关数据~"}
var DATA_CHANGE = ResponseCode{Code: 5, Message: "数据没有任何更改~"}
var DATA_REPEAT = ResponseCode{Code: 6, Message: "数据已存在~"}

//ResponseOk 返回ok
func ResponseOk(data interface{}) *BaseResVO {
	result := new(BaseResVO)
	result.Code = 200
	result.Message = "success"
	result.Data = data
	return result
}

//ResponseError 返回错误
func ResponseError(res *ResponseCode) *BaseResVO {
	result := new(BaseResVO)
	result.Code = res.Code
	result.Message = res.Message
	//数据为空
	result.Data = make(map[string]interface{})
	return result
}
func ResponseErrorCode(code int, message string) *BaseResVO {
	result := new(BaseResVO)
	result.Code = code
	result.Message = message
	//数据为空
	result.Data = make(map[string]interface{})
	return result
}
