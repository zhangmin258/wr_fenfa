package utils

import "time"

// 常量配置
const (
	PAGE_SIZE10    = 10 //列表页每页数据量
	PAGE_SIZE15    = 15
	PAGE_SIZE20    = 20
	PAGE_SIZE40    = 40
	PAGE_SIZE500   = 500
	FormatTime     = "15:04:05"            //时间格式
	FormatDate     = "2006-01-02"          //日期格式
	FormatDateTime = "2006-01-02 15:04:05" //完整时间格式
	WITHDRAWDAY    = "Thursday"
)

// 缓存key
const (
	CACHE_KEY_BindCard       = "wr_fenfa_CACHE_KEY_BindCard_" //绑卡
	CACHE_KEY_User_register  = "user_registerwr_fenfa_"       //验证用户的状态5次处理
	CacheKeyUserPrefix       = "wr_fenfa_CacheKeyUserPrefix_"
	CacheKeyUserInfo         = "wr_fenfa_CacheKeyUserInfo_"
	CacheKeySystemLogs       = "wr_fenfa_CacheKeySystemLogs_"
	CacheKeyOutUsersMenuTree = "wr_fenfa_CacheKeyOutUsersMenuTree_"
	CacheKeyRoleMenuTree     = "wr_fenfa_CacheKeyRoleMenuTree_"
	CacheKeyOutUsersMenuList = "wr_fenfa_CacheKeyOutUsersMenuList_"
	CacheKeyRoleMenuList     = "wr_fenfa_CacheKeyRoleMenuList_"
	CacheKeySystemMenu       = "wr_fenfa_CacheKeySystemMenu" //菜单key
)

// 缓存时间
const (
	RedisCacheTime_User         = 15 * time.Minute
	RedisCacheTime_TwoHour      = 2 * time.Hour
	RedisCacheTime_Role         = 15 * time.Second
	RedisCacheTime_Organization = 24 * time.Hour //24 * time.Hour //组织架构信息缓存时间
)

//信使短信用户密码
const (
	XS_URL      = "http://60.205.151.174:8888/v2sms.aspx"
	XS_NAME     = "微融"
	XS_PASSWORD = "weirong@2017"
	XS_USERID   = "229"
)

// send to users
const (
	ToUsers = "chenxn@zcmlc.com;angzx@zcmlc.com;jgl@zcmlc.com"
	//ToUsers = "chenxn@zcmlc.com"
)

//验证配置
const (
	Regular = "^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0-9])|(17[0-9]))\\d{8}$"
)

const (
	SOURCE              = "weirong"
	CACHE_KEY_Vcode     = "vcode_wr_fenfa_"     //发送验证码
	CACHE_KEY_Pay_Vcode = "pay_vcode_fenfa_wr_" //修改支付密码验证码
)

//连连银行卡
const (
	CACHE_KEY_All_BANKCARD_MAP        = "wr_fenfa_CACHE_KEY_All_BANKCARD_MAP" //支持的银行map
	CACHE_KEY_ChECKWITHDRAWCASH_DEPID = "wr_fenfa_CACHE_KEY_withdrawdeposit_depid_"
	CACHE_KEY_CheckCardNumber         = "wr_fenfa_CACHE_KEY_CheckCardNumber_" //验证银行卡号
	CACHE_KEY_ChECKWITHDRAWCASH_UID   = "wr_fenfa_CACHE_KEY_withdrawcash_uid_"
	CACHE_KEY_CHARGE_UID              = "wr_fenfa_CACHE_KEY_charge_uid_"
	CACHE_KEY_LLAPI_TOKEN_            = "wr_fenfa_CACHE_KEY_LLAPI_TOKEN_"            //连连API绑卡token
	CACHE_KEY_Loan_Accept_State_Deal  = "wr_fenfa_CACHE_KEY_Loan_Accept_State_Deal_" //微融分发放款处理
)

var (
	LL_oid_partner     string                                                               //连连商户号
	LL_RSA_PRIVATE_KEY string                                                               //连连RSA私钥
	LL_RSA_PUBLIC_KEY  string                                                               //连连rsa公钥
	LL_PUBLIC_KEY      string                                                               //连连rsa公钥
	LL_Trade_Query     = "https://instantpay.lianlianpay.com/paymentapi/queryPayment.htm"   // 连连放款订单查询
	LL_Trade_API       = "https://instantpay.lianlianpay.com/paymentapi/payment.htm"        //连连实时付款API
	LL_Trade_CONFIRM   = "https://instantpay.lianlianpay.com/paymentapi/confirmPayment.htm" //连连确认付款地址
	LL_Trade_CallBack  = "notify/lltrade"                                                   //连连放款回调地址

)

const (
	CACHE_KEY_Ip_REGISTER = "wr_fenfa_Operation_Ip_Register_"
	CACHE_KEY_Ip_FORGOT   = "wr_fenfa_Operation_Ip_Forgot_"
)

const (
	WithdrawCount = 10
	WithdrawMoney = 100 // 最低提现金额
)

//新浪长链接转短链接
const (
	LANGTOSHORTURL = "https://api.weibo.com/2/short_url/shorten.json"
	ACCESSTOKEN    = "2.00Lc7rvGH3sUwC0440fd10cblZoJGE"
)
