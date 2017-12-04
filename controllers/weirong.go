package controllers

import (
	"github.com/astaxie/beego"
	"wr_fenfa/models"
	"wr_fenfa/services"
	"wr_fenfa/utils"
	"strconv"
	"wr_fenfa/cache"
	"zcm_tools/email"
)

/*
	提供给微融其他项目的接口
*/
type WeirongController struct {
	beego.Controller
}

// 不校验跨域
func (c *WeirongController) Prepare() {
	c.EnableXSRF = false
}

// 当有用户注册代理的外放链接时，微融将该新用户信息推送至该接口
// 这里我们将新用户作为非代理用户插入到用户表，并与其上级代理添加代理关系
func (c *WeirongController) PostFenfaAddUsers() {

	var err error
	smsStr := ""
	router := "weirong/postfenfaaddusers"
	defer func() {
		if err != nil {
			email.Send(smsStr, c.Ctx.Input.IP()+router+";err:"+err.Error(), utils.ToUsers,"weirong" )
		}
		c.ServeJSON()
	}()
	if c.Ctx.Input.IP() != "127.0.0.1" && c.Ctx.Input.IP() != "192.168.1.233" {
		smsStr = "非法请求！"
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "非法请求！"}
		return
	}
	account := c.GetString("Account")
	fenfaCode := c.GetString("FenfaCode")
	fenfaCode = string(utils.DesBase64Decrypt([]byte(fenfaCode)))
	if account == "" || fenfaCode == "" {
		smsStr = "请求参数异常！"
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "请求参数异常！"}
		return
	}
	// 判断该用户是否存在
	count, err := models.CheckIsRegiste(account)
	if err != nil {
		beego.Debug(err.Error())
	}
	if count == 0 {
		// 向用户表里添加该用户记录（为普通用户）并添加代理关系
		err = services.InsertSimpleUser(fenfaCode, account)
		if err != nil {
			smsStr = "向用户表里添加该用户记录异常！"
			cache.RecordLogs(0, 0, "", "", "", "向用户表里添加该用户记录异常", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"err": "向用户表里添加该用户记录异常!", "ret": 403}
			return
		}
	}

	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}

// 当用户在微融完成第三方产品的注册动作时，将产品id和用户信息传递过来
// 将信息保存
func (c *WeirongController) SendUsersRegister() {
	var err error
	smsStr := ""
	router := "weirong/sendusersregister"
	defer func() {
		if err != nil {
			email.Send(smsStr, c.Ctx.Input.IP()+router+";err:"+err.Error(), utils.ToUsers,"weirong" )
		}
		c.ServeJSON()
	}()
	if c.Ctx.Input.IP() != "127.0.0.1" {
		smsStr = "非法请求!"
		cache.RecordLogs(0, 0, "", "", c.Ctx.Input.IP(), "非法请求", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "非法请求！"}
		return
	}
	account := c.GetString("Account")
	pid := c.GetString("ProductId")
	if account == "" || pid == "" {
		smsStr = "请求参数异常!"
		cache.RecordLogs(0, 0, "", "", "", "请求参数异常", "account||pid为空", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "请求参数异常！"}
		return
	}
	proId, _ := strconv.Atoi(pid)
	count, err := models.IsFenfaUser(account)
	if count < 1 {
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	uid, err := models.GetUidByAccount(account) // 根据account查询用户的uid
	if err != nil {
		smsStr = "根据account查询用户的uid异常!"
		cache.RecordLogs(0, 0, "", "", "", "根据account查询用户的uid异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"err": "根据account查询用户的uid异常", "ret": 403}
		return
	}
	r1, err := services.SearchAllParent(uid) // 获取代理链
	if err != nil {
		smsStr = "获取代理链异常!"
		cache.RecordLogs(0, 0, "", "", "", "获取代理链异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"err": "获取代理链异常!", "ret": 403}
		return
	}
	r2 := utils.ExchangeList(r1)                 // 反序
	p, err := models.GetProductAgentPrice(proId) // 获取产品一级代理价格
	if err != nil {
		smsStr = "获取产品一级代理价格异常!"
		cache.RecordLogs(0, 0, "", "", "", "获取产品一级代理价格异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"err": "获取产品一级代理价格异常!", "ret": 403}
		return
	}
	r3, err := services.SearchAllIncome(r2, p.AgentPrice) // 查询每个代理对应的收益
	if err != nil {
		smsStr = "查询每个代理对应的收益异常!"
		cache.RecordLogs(0, 0, "", "", "", "查询每个代理对应的收益异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"err": "查询每个代理对应的收益异常!", "ret": 403}
		return
	}
	err = services.InsertRegisterIncome(r3, uid, proId, p.Name) // 给每一级代理天添加收益和收益记录
	if err != nil {
		smsStr = "保存每一级代理的收益异常!"
		cache.RecordLogs(0, 0, "", "", "", "保存每一级代理的收益异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"err": "保存每一级代理的收益异常!", "ret": 403}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}
