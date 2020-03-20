package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/songhq211949/beego-api/models"
)

type WuhanController struct {
	beego.Controller
}

func (this *WuhanController) Lists() {
	typeStr := this.GetString("type")
	fmt.Println("前端传入放入参数type为", typeStr)
	o := orm.NewOrm()
	var wuhans []models.Wuhan
	_, err := o.QueryTable("wuhan").All(&wuhans)
	if err != nil {
		fmt.Println("select wuhan is err", err)
		return
	}
	this.Data["json"] = models.ResponseOk(wuhans)
	this.ServeJSON()
}
