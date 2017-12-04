package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/cache"
	"wr_fenfa/utils"
	"wr_fenfa/services"
	"zcm_tools/email"
)

// 用户接口
type UserController struct {
	BaseController
}

func (c *UserController) UserInfo() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	user := c.User
	// 获取用户详细信息
	userInfo, err := models.GetUserInfo(user.Id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取用户详细信息失败", err.Error(), c.Ctx.Input)
	}
	// 用户专属链接
	userInfo.Url = utils.UserUrl + string(utils.DesBase64Encrypt([]byte(user.Code)))
	// 长链接转短链接
	userInfo.Url, err = services.LongToSort(userInfo.Url)
	if err != nil {
		email.Send("长链接转短链接异常!",c.Ctx.Input.IP()+" err: "+err.Error(), utils.ToUsers, "weirong")
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "长链接转短链接异常", err.Error(), c.Ctx.Input)
	}
	c.Data["userInfo"] = userInfo
	c.TplName = "personal-center/personal-info.html"
}
