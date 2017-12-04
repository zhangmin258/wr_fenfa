package models

import (
	"github.com/astaxie/beego/orm"
)

// 绑卡
type BindCardRequest struct {
	User
	BankId         int    // 银行id
	BankName       string // 银行名称
	BankCardNumber string // 银行卡号
	Channel        int    // 通道
	Data           string
}

// 用户提交的绑卡信息
type CardInfo struct {
	User
	UserName       string // 用户实名
	IdCard         string // 身份证号
	BankMobile     string // 银行预留手机号
	BankId         int    // 银行id
	BankName       string // 银行名称
	BankCardNumber string // 银行卡号
}

type Card struct {
	Id             int
	BankCardNumber string // 银行卡号
}

// 连连绑卡返回信息
type BindCardData struct {
	Ret_code    string
	Ret_msg     string
	Sign_type   string
	Sign        string
	Oid_partner string
	User_id     string
	Result_sign string
	No_agree    string
}

// 聚合绑卡四要素检测结果
type JuheCheckBankCardResult struct {
	Reason string `json:"reason"`
	Result struct {
		Jobid    string `json:"jobid"`
		Realname string `json:"realname"`
		Bankcard string `json:"bankcard"`
		Idcard   string `json:"idcard"`
		Mobile   string `json:"mobile"`
		Res      string `json:"res"`
		Message  string `json:"message"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

// 判断银行卡号是否已经存在--再使用中
func IsExistBankCardNumber(cardnum string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT count(1) FROM users_bankcards WHERE card_number = ?`
	err = o.Raw(sql, cardnum).QueryRow(&count)
	return count, err
}

// 用户是否绑卡
func IsBandingCard(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_bankcards WHERE uid=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

// 保存用户绑卡信息
func SaveUserBankCard(m CardInfo) error {
	sql := `INSERT INTO users_bankcards (uid,real_name,card_number,id_card,bank_mobile,create_time) VALUES(?,?,?,?,?,NOW())`
	_, err := orm.NewOrm().Raw(sql, m.Id, m.UserName, m.BankCardNumber, m.IdCard, m.BankMobile).Exec()
	return err
}

// 根据id获取银行卡信息
func GetBankCardById(id int) (card Card, err error) {
	sql := `SELECT id,card_number AS bank_card_number FROM users_bankcards WHERE uid = ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&card)
	return
}
