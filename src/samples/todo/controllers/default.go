package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	beego.Info("MainController --------------------------------------------------------")
	this.TplName = "index.html"
	this.Render()
}
