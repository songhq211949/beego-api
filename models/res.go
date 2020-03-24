package models

const LOGIN   = 1
const PING  = 0

//UserLoginDTO  用户登入
type UserLoginDTO struct{
	Uid int `json:"uid"`
}
