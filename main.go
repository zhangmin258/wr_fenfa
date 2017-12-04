package main

import (
	_ "wr_fenfa/routers"
	"github.com/astaxie/beego"
	"wr_fenfa/services"
	"wr_fenfa/utils"
)

func main() {
	services.Task()
	beego.Run()
}

func init()  {
	beego.AddFuncMap("accountDispose", utils.AccountDispose)
	beego.AddFuncMap("formatTimeToString", utils.FormatTimeToString)
	beego.AddFuncMap("bankCardFormat", utils.BankCardFormat)
	beego.AddFuncMap("float64ToString", utils.Float64ToString)
}

