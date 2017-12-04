package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"time"
	"wr_fenfa/services"
	"encoding/json"
	"wr_fenfa/cache"
)

// 用户绑卡接口
type UsersBankCardController struct {
	BaseController
}

// 绑卡页面
func (c *UsersBankCardController) BindCardPage() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.Data["user"] = c.User
	c.TplName = "bind-card.html"
}

// 用户绑卡接口
func (c *UsersBankCardController) BindCard() {
	var m models.CardInfo
	m.Id = c.User.Id
	var errmsg string
	defer func() {
		if errmsg != "" {
			utils.Rc.Delete(utils.CACHE_KEY_CheckCardNumber + m.BankCardNumber)
		}
		c.ServeJSON()
	}()
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &m)
	if err != nil {
		errmsg = "解析用户参数出错"
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "解析用户参数出错", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "解析用户参数出错"}
		return
	}
	if !utils.Rc.SetNX(utils.CACHE_KEY_CheckCardNumber+m.BankCardNumber, 1, time.Second*30) {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "重复提交","", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "亲,你已提交请求，请稍后再试~"}
		return
	}
	if m.BankCardNumber == "" {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "参数有误，请稍后再试~", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "参数有误，请稍后再试~"}
		errmsg = "BankId/BankCardNumber/BankName为空"
		return
	}
	count, err := models.IsExistBankCardNumber(m.BankCardNumber)
	if err != nil {
		errmsg = "判断卡号是否已经存在出错" + err.Error()
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "判断卡号是否已经存在出错", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "系统错误，请稍后再试~"}
		return
	}
	if count > 0 {
		errmsg = "卡号已存在"
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "卡号已存在", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "卡号已存在~"}
		return
	}

	// 判断改用户是否有绑卡
	count1, err := models.IsBandingCard(m.Id)
	if err != nil {
		errmsg = "判断卡号是否已经存在出错" + err.Error()
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "判断卡号是否已经存在出错", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "系统错误，请稍后再试~"}
		return
	}
	if count1 > 0 {
		errmsg = "该用户已绑卡"
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "该用户已绑卡", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "您已经绑定银行卡~"}
		return
	}
	// 调用聚合认证四元素
	res, err := services.JuheCheckBankCard(m.UserName, m.BankCardNumber, m.IdCard, m.BankMobile)
	if err != nil {
		errmsg = "调用聚合认证异常"
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "调用聚合认证异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "系统错误，请稍后再试~"}
		return
	}
	if res.ErrorCode == 0 && res.Result.Res == "1" { // 验证成功
		// 保存用户银行卡信息
		err = models.SaveUserBankCard(m)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "保存用户银行卡错误", err.Error(), c.Ctx.Input)
			errmsg = "保存用户银行卡错误：" + err.Error()
			c.Data["json"] = map[string]interface{}{"ret": 403, "err": "保存用户银行卡错误"}
			return
		}
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	} else {
		errmsg = "输入信息核对错误！"
		cache.RecordLogs(c.User.Id, 0, c.User.Account, "", "", "输入信息核对错误", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "您输入的信息有误，请核对后重试！~"}
		return
	}
}
