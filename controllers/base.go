package controllers

import (
	"github.com/astaxie/beego"
	"net/url"
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"zcm_tools/log"
	"encoding/json"
	"net/http"
	"encoding/base64"
	"wr_fenfa/cache"
	"wr_fenfa/services"
)

var v1Log *log.Log

func init() {
	v1Log = log.Init()
}

// BaseController 基础controller
type BaseController struct {
	beego.Controller
	User *models.User
}

// 前置方法
func (c *BaseController) Prepare() {
	// 请求体解码
	if utils.H5Encoded == "true" {
		// 请求体base64解码
		if c.Ctx.Input.Method() == "POST" && len(c.Ctx.Input.RequestBody) != 0 {
			var err error
			c.Ctx.Input.RequestBody, err = base64.StdEncoding.DecodeString(string(c.Ctx.Input.RequestBody))
			if err != nil {
				c.Data["json"] = map[string]interface{}{"ret": 403, "err": "网络异常！"}
				// TODO 日志
				// cache.RecordLogs(0, "0", "解析请求体失败！"+err.Error(), c.Ctx.Input.IP(), "h5")
				c.ServeJSON()
				c.StopRun()
			}
		}
	}
	// 验证缓存tiket
	verify := false
	userTiket := c.Ctx.GetCookie("weirong_tiket")
	userCode := c.Ctx.GetCookie("weirong_code")
	if userTiket != "" {
		if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyUserPrefix+userCode) { // 判断Code是否存在
			if tiket, err := utils.Rc.RedisString(utils.CacheKeyUserPrefix + userCode); err == nil { // 比对tiket
				if tiket == userTiket {
					u, err := cache.GetUserByCodeCache(userCode) // 查询用户userInfo
					if err != nil {
						c.Data["json"] = map[string]interface{}{"ret": 403, "err": "用户信息异常！"}
						c.ServeJSON()
						c.StopRun()
					}
					c.User = u
					// 菜单权限校验
					if c.Ctx.Input.URL() != "/" {
						b, err := services.CheckMenu(c.Ctx.Input.URL(),u)
						if err != nil {
							beego.Debug(err.Error())
							c.Data["json"] = map[string]interface{}{"ret": 403, "err": "校验访问权限出错！"}
							c.ServeJSON()
							c.StopRun()
						}
						if !b {
							c.Data["json"] = map[string]interface{}{"ret": 403, "err": "访问权限不足！"}
							// 清除缓存
							if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySystemMenu) {
								utils.Rc.Delete(utils.CacheKeySystemMenu)
							}
							if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyOutUsersMenuList) {
								utils.Rc.Delete(utils.CacheKeyOutUsersMenuList)
							}
							c.ServeJSON()
							c.StopRun()
						}
					}
					ip := c.Ctx.Input.IP()
					requestBody, _ := url.QueryUnescape(string(c.Ctx.Input.RequestBody))
					v1Log.Println("请求地址：", c.Ctx.Input.URI(), "用户信息：", requestBody, "RequestBody：", requestBody, "IP：", ip)
					// 更新tiket时间
					if ip == "127.0.0.1" || ip == "60.191.125.34" || ip == "60.191.37.251" {
						utils.Rc.Put(utils.CacheKeyUserPrefix+userCode, userTiket, utils.RedisCacheTime_TwoHour)
					} else {
						utils.Rc.Put(utils.CacheKeyUserPrefix+userCode, userTiket, utils.RedisCacheTime_User)
					}
					verify = true
					// TODO 更新用户最近操作时间
				}
			}
		}
	}
	// 上传文件跳过验证
	if c.Ctx.Input.IsUpload() {
		verify = true
	}
	if !verify {
		if c.Ctx.Input.IsAjax() {
			c.Ctx.Output.JSON(map[string]interface{}{"ret": 408, "msg": "timeout"}, false, false)
			c.StopRun()
		} else {
			c.Ctx.Redirect(302, "/login")
			c.StopRun()
		}
	}
	// 用户没有绑卡则跳转到绑卡页面
	if c.User.AccountType == 0 && c.Ctx.Input.URL() != "/usersbankcard/bindcardpage" && c.Ctx.Input.URL() != "/usersbankcard/bindcard" {
		if count, _ := models.IsBandingCard(c.User.Id); count == 0 {
			c.Ctx.Redirect(302, "/usersbankcard/bindcardpage")
			c.StopRun()
		}
	}
}

// 是否需要模板
func (c *BaseController) IsNeedTemplate() {
	pushstate := c.GetString("pushstate")
	if pushstate != "1" {
		c.Data["DisplayName"] = c.User.DisplayName
		c.Layout = "layout/layout.html"
	}
}

func (c *BaseController) ServeJSON(encoding ...bool) {
	var (
		hasIndent   = true
		hasEncoding = false
	)
	if beego.BConfig.RunMode == beego.PROD {
		hasIndent = false
	}
	if len(encoding) > 0 && encoding[0] == true {
		hasEncoding = true
	}
	c.JSON(c.Data["json"], hasIndent, hasEncoding)
}

func (c *BaseController) JSON(data interface{}, hasIndent bool, coding bool) error {
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(c.Ctx.Output.Context.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding {
		content = []byte(utils.StringsToJSON(string(content)))
	}
	// app接口都是页面接口，不需要处理加密问题。
	// H5接口返回json需要base64编码
	if utils.H5Encoded == "true" {
		return c.Ctx.Output.Body(utils.Base64Encrypt(content))
	}
	return c.Ctx.Output.Body(content)
}
