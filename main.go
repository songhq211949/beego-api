package main

import (
	"github.com/astaxie/beego"
	_ "github.com/songhq211949/beego-api/routers" //这里引入routers模块。会调用init方法
)

func main() {
	beego.Run()
}
