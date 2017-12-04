package controllers

import (
	"wr_fenfa/services"
	"wr_fenfa/utils"
	"fmt"
)

// HomeController 主页
type HomeController struct {
	BaseController
}

// Get 主页Get
func (c *HomeController) Get() {
	xsrfToken := c.XSRFToken()
	fmt.Println(xsrfToken)
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	c.TplName = "index.html"
}

// Post 主页获取数据
func (c *HomeController) Post() {
	m, err := services.GetSysMenuTree(c.User) // 获取系统菜单
	if err != nil && err.Error() != utils.ErrNoRow() {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "data": m}
	}
	c.ServeJSON()
}
