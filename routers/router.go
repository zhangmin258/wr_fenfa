package routers

import (
	"wr_fenfa/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AccountController{}, "get:Login;post:CheckPassword") //登录页
	beego.Router("/loginout", &controllers.AccountController{}, "get:LoginOut")              //退出登录
	beego.MyAutoRouter(&controllers.AccountController{})
	beego.MyAutoRouter(&controllers.SystemController{})
	beego.MyAutoRouter(&controllers.UsersBankCardController{})
	beego.MyAutoRouter(&controllers.DataController{})
	beego.MyAutoRouter(&controllers.StaffController{})
	beego.MyAutoRouter(&controllers.UserController{})
	beego.MyAutoRouter(&controllers.WithdrawController{})
	beego.MyAutoRouter(&controllers.WeirongController{})
	beego.MyAutoRouter(&controllers.CommissionController{})
	beego.MyAutoRouter(&controllers.AdminController{})
	beego.MyAutoRouter(&controllers.NotifyController{})
}
