package models

//重构接口,基本返回对象
type BaseResVO struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOk(data interface{}) *BaseResVO {
	result := new(BaseResVO)
	result.Code = 200
	result.Message = "success"
	result.Data = data
	return result
}
