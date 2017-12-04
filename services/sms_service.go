package services

import (
	"wr_fenfa/utils"
	"time"
	"zcm_tools/email"
	"github.com/astaxie/beego"
	"wr_fenfa/cache"
	"github.com/astaxie/beego/context"
)

type XsSmsResult struct {
	ReturnStatus  string `xml:"returnstatus"`
	Message       string `xml:"message"`
	RemainPoint   string `xml:"remainpoint"`
	TaskId        string `xml:"taskID"`
	SuccessCounts int    `xml:"successCounts"`
}

// 获取注册验证码
func GetVcode(account, source string, input *context.BeegoInput) {
	vcode := utils.GetRandDigit(4) // 获取4位有效的验证码
	content := "【微融】验证码为" + vcode + "，感谢您注册微融分发平台。请在5分钟内完成注册，工作人员不会向您索取验证码，请勿泄露。"
	var result string
	if len(account) == 11 {
		xs, err := SendXinShimessage(content, account)
		if err != nil || xs.ReturnStatus != "Success" {
			result = "验证码发送失败~"
			email.Send("验证码发送失败", "请求参数:"+"account:"+account+";content:"+content+"source:"+source+";方法:GetForgetPwdVcode;err:"+err.Error(), utils.ToUsers, "weirong")
		}
	}
	cache.RecordLogs(0, 0, "", "", "account", "发送用户验证码", "注册验证码为"+vcode+"==========结果："+result, input)
	if err := utils.Rc.Put(utils.CACHE_KEY_Vcode+account, vcode, time.Minute*5); err != nil {
		beego.Error("vcode_wr::" + err.Error())
	}
}

//忘记密码--找回登陆密码验证码
func GetForgetPwdVcode(account, source string, input *context.BeegoInput) {
	vcode := utils.GetRandDigit(4) // 获取4位有效的验证码
	content := "【微融】您的验证码为" + vcode + "，您正在申请找回微融分发平台密码,请在5分钟内完成操作，工作人员不会向您索取验证码，请勿泄露。"
	var result string
	if len(account) == 11 {
		xs, err := SendXinShimessage(content, account)
		if err != nil || xs.ReturnStatus != "Success" {
			result = "验证码发送失败~"
			email.Send("验证码发送失败", "请求参数:"+"account:"+account+";content:"+content+"source:"+source+";方法:GetForgetPwdVcode;err:"+err.Error(), utils.ToUsers, "weirong")
		}
	}
	cache.RecordLogs(0, 0, "", "", "account", "发送用户验证码", "注册验证码为"+vcode+"==========结果："+result, input)
	if err := utils.Rc.Put(utils.CACHE_KEY_Vcode+account, vcode, 5*time.Minute); err != nil {
		beego.Error(utils.CACHE_KEY_Vcode + "::" + err.Error())
	}
}
