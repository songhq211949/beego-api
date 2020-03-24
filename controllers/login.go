package controllers

import (
	"github.com/astaxie/beego"
	"github.com/songhq211949/beego-api/models"
)

type LoginController struct{
	beego.Controller
}

func (c *LoginController) LoginError(){
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		c.StopRun()
	}