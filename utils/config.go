package utils

import (
	"zcm_tools/cache"
	"github.com/astaxie/beego"
)

// MYSQL
var (
	RunMode          string // 运行模式
	MYSQL_URL        string // 主库
	MYSQL_LOG_URL    string // log库
	MYSQL_BACKUP_URL string // 备份库
	MYSQL_FENFA_URL  string // 分发平台库
)

const (
	WR_FENFA_API_URL ="https://ff.weiyunjinrong.cn/" // 微融分发后台地址
)

// Redis
var (
	Rc          *cache.Cache // redis缓存
	Re          error        // redis错误
	BEEGO_CACHE string       // redis地址
)

// base
var (
	Enablexsrf string // XSRF校验开关
	H5Encoded  string // H5接口base64编码开关
)

var (
	UserUrl string // 用户的专属链接
)

func init() {
	RunMode = beego.AppConfig.String("run_mode")
	config, err := beego.AppConfig.GetSection(RunMode)
	if err != nil {
		panic("配置文件读取错误 " + err.Error())
	}
	Enablexsrf = beego.AppConfig.String("enablexsrf")
	H5Encoded = beego.AppConfig.String("h5_encoded")
	// mysql
	MYSQL_URL = config["mysql_url"]
	MYSQL_LOG_URL = config["mysql_log_url"]
	MYSQL_BACKUP_URL = config["mysql_backup_url"]
	MYSQL_FENFA_URL = config["mysql_fenfa_url"]
	// redis
	BEEGO_CACHE = config["beego_cache"]
	Rc, Re = cache.NewCache(BEEGO_CACHE)
	// show
	beego.Info("┌───────────────────")
	beego.Info("│模式:" + RunMode)
	beego.Info("│XSRF校验:" + Enablexsrf)
	beego.Info("│H5接口编码:" + H5Encoded)
	beego.Info("└───────────────────")
	if RunMode == "release" {
		UserUrl = "https://weiyunjinrong.cn/user/invitationfriends?source=fenfa%26ucode="
	} else {
		UserUrl = "http://192.168.1.233:7007/user/invitationfriends?source=fenfa%26ucode="
	}
	//连连快捷支付
	LL_oid_partner = "201710010000979003"
	LL_RSA_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBALvSN3AQKs0NIGbW
C0SL/aRiH+OF+IgLPWHFs3kMEnpW1+WjqCoJE9Nw99NNkZKEW552g8ejh3GISJYJ
3njNfqMCflKgsRzcncrZR3DSVgQM9GY5HLqDLzkAIrJ/EXCBfS/tbKzX8zJ+X5W/
MGYspn+w10YrqhjCV3wPO8kFNddxAgMBAAECgYBHsho5a+J6va0FtGU+uFWNP2u+
1XAmtmuq++XjqikPjEEDxvI1gZuQ1gm0HmMYU/AJUGJDffgA7a4PoBrNcFwLREfC
xqtrqnAfa5Ub+Xat/KPtd8hdQnC6JrEcpeLCKAjWsVLRAm/4/iinpa2xcueBv0Og
aC2iOYSRwotObdhPhQJBAN3poj8hf7l0JytBAviBmKHNOBLO4wY1RAyh2aP4bWnu
VZaJ5VnHra47ZrD2VqfvGyN2EfaprEtpPYEwtWodovcCQQDYq/0yNfAwrsMwufEQ
NjETpBVrWR49ztbHpuu3IneL0M7BAezcf5QuhjfkvRKsoC11vNBAr/enIWBC5r/j
RVbXAkBSRXDyaM/6iHaREawxR5K3weadCnieb5cH++U9ZjfiQwsWIY+XJnFcnAcp
alqcLghosDhes27+EklMISvQ6KXnAkBtEXatNdWoy/BZsOAWRxFBT9GwbfX5KwuX
CQGS+HixGvVY1v1Cqb4QBWRRcpPZ7e+0Ws2CIpJJwVVRmBJz902VAkBHmNtAe05P
xiV6/vxS53R12vVdaagCCt8sy/m0yMxZUVvzQoVMxqn3ABaG5qYaX+cNQjJ4CRI8
3IxpfMg5jWzP
-----END PRIVATE KEY-----
`
	LL_RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC70jdwECrNDSBm1gtEi/2kYh/j
hfiICz1hxbN5DBJ6Vtflo6gqCRPTcPfTTZGShFuedoPHo4dxiEiWCd54zX6jAn5S
oLEc3J3K2Udw0lYEDPRmORy6gy85ACKyfxFwgX0v7Wys1/Myfl+VvzBmLKZ/sNdG
K6oYwld8DzvJBTXXcQIDAQAB
-----END PUBLIC KEY-----
`

	LL_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSS/DiwdCf/aZsxxcacDnooGph3d2JOj5GXWi+
q3gznZauZjkNP8SKl3J2liP0O6rU/Y/29+IUe+GTMhMOFJuZm1htAtKiu5ekW0GlBMWxf4FPkYlQ
kPE0FtaoMP3gYfh+OwI+fIRrpW3ySn3mScnc6Z700nU/VYrRkfcSCbSnRwIDAQAB
-----END PUBLIC KEY-----
`
}
