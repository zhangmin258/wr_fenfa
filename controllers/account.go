package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"github.com/astaxie/beego"
	"zcm_tools/uuid"
	"time"
	"encoding/json"
	"wr_fenfa/cache"
	"zcm_tools/email"
	"strings"
	"wr_fenfa/services"
)

// 登录控制器
type AccountController struct {
	beego.Controller
}

// 登录页面
func (c *AccountController) Login() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.TplName = "login.html"
}

// 登陆业务
func (c *AccountController) CheckPassword() {
	defer c.ServeJSON()
	account := c.GetString("account")
	password := c.GetString("password")
	// 根据手机号获取用户code
	code, err := models.GetCodeByAccount(account)
	if err != nil && err.Error() == utils.ErrNoRow() {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": "登录失败！该账号尚未注册"}
		return
	}
	// 验证登陆
	password = utils.MD5(code + password)
	m, _ := models.Login(account, password)
	if m == nil {
		c.Data["json"] = map[string]interface{}{"ret": 404, "err": "登录失败！用户名或密码不正确"}
		return
	} else if m.IsUsed == 1 {
		c.Data["json"] = map[string]interface{}{"ret": 404, "err": "账号已禁用，请联系客服或上级代理!"}
		return
	}
	// 生成tiket
	tiket := "wr_fenfa" + time.Now().Format("20060102") + uuid.NewUUID().Hex()
	// 存入缓存
	if utils.Re != nil {
		c.Data["json"] = map[string]interface{}{"ret": 404, "err": "登录失败！请稍后再试！"}
		return
	}
	data, err := json.Marshal(m)
	if err != nil {
		cache.RecordLogs(0, 0, account, "", "Login", "序列化用户数据失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 404, "err": "登录失败！用户信息异常！"}
		return
	}
	utils.Rc.Put(utils.CacheKeyUserPrefix+m.Code, tiket, utils.RedisCacheTime_User) // 用户的登陆凭证
	utils.Rc.Put(utils.CacheKeyUserInfo+m.Code, data, utils.RedisCacheTime_User)    // 用户的详情
	// 存入cookie
	c.Ctx.SetCookie("weirong_code", m.Code)
	c.Ctx.SetCookie("weirong_tiket", tiket)
	var banding = 0  // 外部用户判断是否绑卡 0:已绑卡，1：未绑卡
	var userType = 1 // 0 外部用户 1 内部用户
	if m.AccountType == 0 {
		userType = 0
		if count, _ := models.IsBandingCard(m.Id); count == 0 {
			banding = 1
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "登录成功", "jumpBand": banding, "userType": userType}
	return
}

// 退出登录
func (c *AccountController) LoginOut() {
	weirong_code := c.Ctx.GetCookie("weirong_code")
	user, err := models.GetUserInfoByCode(weirong_code)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "Login", "根据code查询用户失败，code:"+weirong_code, "", c.Ctx.Input)
	}
	if weirong_code != "" {
		cache.RecordLogs(user.Id, 0, user.Account, "", "Login", "记录用户下线", "用户code："+weirong_code, c.Ctx.Input)
		// 清除cookie
		c.Ctx.SetCookie("weirong_code", "", -1)
		c.Ctx.SetCookie("weirong_tiket", "", -1)
		// 清除redis
		utils.Rc.Delete(utils.CacheKeyUserPrefix + weirong_code)
	}
	c.Ctx.Redirect(302, "/login")
}

// 忘记密码
func (c *AccountController) ForgotPasswordPage() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.TplName = "update-password.html"
}

// 外部用户注册
func (c *AccountController) RegisterPage() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.TplName = "reg.html"
}

// 注册
func (c *AccountController) Register() {
	var user models.UserRequest
	err := c.ParseForm(&user)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "", "解析用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析用户失败"}
		return
	}
	defer func() {
		utils.Rc.Delete("Registerwrfenfa_" + user.Account)
		c.ServeJSON()
	}()
	// 账号为空
	if user.Account == "" {
		c.Data["json"] = map[string]interface{}{"err": "请填写账号!", "ret": 403}
		return
	}
	// 验证手机号
	if !utils.Validate(user.Account) {
		c.Data["json"] = map[string]interface{}{"err": "手机号码错误!", "ret": 403}
		return
	}
	// 判断该用户是否已为代理
	count, err := models.CheckAgentUser(user.Account)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "判断用户出错!"}
		return
	}
	if count > 0 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "该用户已存在!"}
		return
	}
	if user.Vcode == "" {
		c.Data["json"] = map[string]interface{}{"err": "请填写验证码!", "ret": 403}
		return
	}
	// 判断验证码是否过期
	key := utils.CACHE_KEY_Vcode + user.Account
	code, err := utils.Rc.RedisBytes(key)
	if err != nil || !utils.Rc.IsExist(key) {
		c.Data["json"] = map[string]interface{}{"err": "验证码过期，请重新发送！", "ret": 403}
		if err != nil {
			cache.RecordLogs(0, 0, "", "", "", "验证码错误,验证码过期，请重新发送！", err.Error(), c.Ctx.Input)
		}
		return
	}
	// 判断验证码
	if !strings.EqualFold(string(code), user.Vcode) {
		c.Data["json"] = map[string]interface{}{"err": "验证码错误！", "ret": 403}
		return
	}
	// 控制频繁提交
	if !utils.Rc.SetNX("Registerwrfenfa_"+user.Account, 1, time.Minute) {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "亲,你已提交请求，请稍后再试~"}
		return
	}
	// 判断用户是否已为普通用户
	count, _ = models.CheckUser(user.Account)
	var uid int
	if count > 0 {
		uid, err = models.GetUidByAccount(user.Account)
		if err != nil {
			beego.Debug(err.Error())
		}
		code, _ := models.GetCodeByAccount(user.Account)
		user.Password = utils.MD5(code + user.Password)
		// 普通用户，将其升级为代理
		err = models.UpdateSimpleUsersBeAgent(&user.User)
		if err != nil {
			email.Send("用户注册时将普通用户升级为代理出错", "err:"+err.Error(), utils.ToUsers, "weirong")
			c.Data["json"] = map[string]interface{}{"err": "升级普通用户为代理失败!", "ret": 403}
			cache.RecordLogs(0, 0, "", "", "", "添加注册用户失败", err.Error(), c.Ctx.Input)
			return
		}
	} else {
		// 直接代理，新增加一条代理信息
		user.Code = utils.NewUUID().Hex()
		user.Password = utils.MD5(user.Code + user.Password)
		user.AccountType = 0
		user.DisplayName = user.Account
		user.IsUsed = 0
		user.AgentLevel = 1
		user.DisplayName = user.Account
		user.IsAgent = 1
		uid, err = models.AddUsers(&user.User)
		if err != nil {
			email.Send("用户注册时添加注册用户出错", "err:"+err.Error(), utils.ToUsers, "weirong")
			c.Data["json"] = map[string]interface{}{"err": "添加注册用户失败!", "ret": 403}
			cache.RecordLogs(0, 0, "", "", "", "添加注册用户失败", err.Error(), c.Ctx.Input)
			return
		}
	}
	// 初始化钱包
	err = models.InnitWallet(uid)
	if err != nil {
		email.Send("初始化钱包出错", "err:"+err.Error(), utils.ToUsers, "weirong")
		c.Data["json"] = map[string]interface{}{"err": "初始化钱包失败!", "ret": 403}
		cache.RecordLogs(0, 0, "", "", "", "初始化钱包失败", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "注册成功"}
}

// 查询帐号注册的时候发送验证码
// @router /existaccount [post]
func (c *AccountController) IsExistAccount() {
	defer c.ServeJSON()
	account := c.GetString("phone")
	if !utils.Validate(account) {
		c.Data["json"] = map[string]interface{}{"ret": 500, "err": "手机号码错误!"}
		return
	}
	// 判断是否有该用户
	count, err := models.CheckAgentUser(account)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "", "判断账号是否存在时通过手机号查询用户信息失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "查询用户信息失败"}
		return
	}
	// 控制频繁提交
	{
		lock := "RegisterFenfaPwdSendVcode_" + account
		if utils.Rc.IsExist(lock) {
			cache.RecordLogs(0, 0, "", "", "", "你已经提交请求，请勿频繁提交!", "", c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 403, "err": "你已经提交请求，请勿频繁提交!"}
			return
		} else {
			utils.Rc.Put(lock, 1, 1*time.Minute)
		}
	}
	// ip校验
	ip := c.Ctx.Input.IP()
	count1 := utils.CheckPwd5Time(utils.CACHE_KEY_Ip_REGISTER + ip)
	if count1 < 0 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "今天的验证码次数已经使用完！"}
		return
	}
	// 查询是否有该用户
	if count == 0 {
		// 用户校验
		count := utils.CheckPwd5Time(utils.CACHE_KEY_User_register + account)
		if count < 0 {
			cache.RecordLogs(0, 0, "", "", "", "今天的验证码次数已经使用完。请明天操作！", "", c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 403, "err": "今天的验证码次数已经使用完。请明天操作！"}
			return
		}
		go services.GetVcode(account, utils.SOURCE, c.Ctx.Input) //调用协程方法,进行查询数据
		//================================================================
		c.Data["json"] = map[string]int{"ret": 200} //未注册返回
		return
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "该号码已经注册哦!"}
		return
	}
}

// 找回密码
func (c *AccountController) ForgotPassword() {
	var user models.UserRequest
	err := c.ParseForm(&user)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "", "用户修改密码时时解析请求参数出错", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析用户失败"}
		return
	}
	defer func() {
		utils.Rc.Delete("Updatewrfenfa_" + user.Account)
		c.ServeJSON()
	}()
	if user.Account == "" {
		c.Data["json"] = map[string]interface{}{"err": "请填写账号!", "ret": 403}
		return
	}
	if !utils.Validate(user.Account) {
		c.Data["json"] = map[string]interface{}{"err": "手机号码错误!", "ret": 403}
		return
	}
	// 判断是否有该用户
	count, err := models.CheckUser(user.Account)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "判断用户出错!"}
		return
	}
	if count == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "该用户不存在!"}
		return
	}
	if user.Vcode == "" {
		c.Data["json"] = map[string]interface{}{"err": "请填写验证码!", "ret": 403}
		return
	}
	// 控制频繁提交
	if !utils.Rc.SetNX("Updatewrfenfa_"+user.Account, 1, time.Minute) {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "亲,你已提交请求，请稍后再试~"}
		return
	}

	// 获取用户code
	user.Code, err = models.GetCodeByAccount(user.Account)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"err": "获取用户code失败", "ret": 403}
		return
	}

	// type控制用户的操作，0是忘记密码，1是修改密码
	if user.Type == 0 { // when type=0,忘记密码，找回密码
		key := utils.CACHE_KEY_Vcode + user.Account
		code, err := utils.Rc.RedisBytes(key)
		if err != nil || !utils.Rc.IsExist(key) {
			c.Data["json"] = map[string]interface{}{"err": "验证码过期，请重新发送！", "ret": 403}
			cache.RecordLogs(0, 0, "", "", "", "验证码错误,验证码过期，请重新发送！", err.Error(), c.Ctx.Input)
			return
		}
		if !strings.EqualFold(string(code), user.Vcode) {
			cache.RecordLogs(0, 0, "", "", "", "验证码错误", "", c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"err": "验证码错误！", "ret": 304}
			return
		}
		// 更新用户信息
		err = models.UpdateUser(user.User)
		if err != nil {
			cache.RecordLogs(0, 0, "", "", "", "设置登录密码失败", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "err": "设置登录密码失败！"}
		}
		c.Data["json"] = map[string]interface{}{"ret": 200}
	} else if user.Type == 1 { // when type=1,有原密码的情况下修改密码

	}
}

// 忘记密码发送验证码
// @router /forgetpwdsendvcode [post]
func (c *AccountController) ForgetPwdSendVcode() {
	defer c.ServeJSON()
	account := c.GetString("phone")
	if !utils.Validate(account) {
		c.Data["json"] = map[string]interface{}{"ret": 500, "err": "手机号码错误!"}
		return
	}
	// 判断是否有该用户
	count, err := models.CheckUser(account)
	if err != nil {
		cache.RecordLogs(0, 0, "", "", "", "判断账号是否存在时通过手机号查询用户信息失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "查询用户信息失败"}
		return
	}
	if count == 0 {
		cache.RecordLogs(0, 0, "", "", "", "用户不存在", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "此用户尚未注册!"}
		return
	}
	//控制频繁提交
	{
		lock := "ForgetFenfaPwdSendVcode_" + account
		if utils.Rc.IsExist(lock) {
			cache.RecordLogs(0, 0, "", "", "", "你已经提交请求，请勿频繁提交!", "", c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 403, "err": "你已经提交请求，请勿频繁提交!"}
			return
		} else {
			utils.Rc.Put(lock, 1, 1*time.Minute)
		}
	}
	// ip校验
	ip := c.Ctx.Input.IP()
	count2 := utils.CheckPwd5Time(utils.CACHE_KEY_Ip_FORGOT + ip)
	if count2 < 0 {
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "今天的验证码次数已经使用完！"}
		return
	}
	//验证码5次限制，递推，会出现10后 中断用户操作
	count1 := utils.CheckPwd5Time(utils.CACHE_KEY_Vcode + account)
	if count1 <= 0 {
		cache.RecordLogs(0, 0, "", "", "", "今天的短信验证码次数已经使用完;请明天操作！", "", c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "今天的短信验证码次数已经使用完;请明天操作！"}
		return
	}
	go services.GetForgetPwdVcode(account, utils.SOURCE, c.Ctx.Input) //调用协程方法,进行查询数据
	//================================================================
	c.Data["json"] = map[string]int{"ret": 200}
	return
}
