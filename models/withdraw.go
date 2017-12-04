package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 用户提现订单信息
type UserWithdraw struct {
	Id             int
	Uid            int       // 用户id
	Account        string    // 用户手机号
	OrderCode      string    // 商户编号
	Amount         float64   // 用户提现金额
	Real_amount    float64   // 实际提现金额
	AmountTime     time.Time // 提现时间
	OrderState     int       // 1.等待处理 2.同意提现 3拒绝提现
	ResultPay      string    // 提现订单处理结果
	CheckPeople    int       // 审批人
	CheckTime      time.Time // 审批时间
	RetCode        string    // 连连提交状态码
	IsIllegalOrder int       // 是否是异常订单0不是 1是
	Finishtime     time.Time // 完成时间
	SearchTime     int       // 主动查询次数
}

// 用户提现订单信息
type UserWithdrawOrder struct {
	Uid        int       // 用户id
	OrderCode  string    // 商户编号
	AmountTime time.Time // 提现时间
	ResultPay  string    // 提现订单处理结果
	RetCode    string    // 连连提交状态码
}

// 用户的绑卡信息
type UsersBankcard struct {
	Id         int
	Uid        int
	RealName   string // 用户真实姓名
	CardNumber string // 银行卡号
	IdCard     string // 身份证
	BankMobile string // 银行预留手机号
	CreateTime time.Time
}

type Pay struct {
	Uid        int
	BankCardId int     //银行卡ID
	Money      float64 //提现金额
	PayPwd     string  // 支付密码
}

// 提现操作参数
type DepositStr struct {
	DepId string // 提现表id
}

// 获取所有佣金列表
func GetUserDeposit(condition string, params []string, begin, size int) (userDeposit []*UserWithdraw, err error) {
	sql := `SELECT ur.id,ur.uid,ur.amount_time,ur.order_state,u.account,ur.amount,ur.result_pay FROM users_withdraw_records ur
	INNER JOIN users u ON ur.uid = u.id
	 LEFT JOIN users_wallet w
	 ON ur.uid=w.uid
	 WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY ur.amount_time DESC LIMIT ?,? "
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&userDeposit)
	return
}

// 获取所有佣金列表
func GetUserDepositCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_withdraw_records ur
	INNER JOIN users u ON ur.uid = u.id
	 LEFT JOIN users_wallet w
	 ON ur.uid=w.uid
		WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 获取用户详细佣金信息
func GetUserDepositDetail(uid int, begin, size int) (userDeposit []*UserWithdraw, err error) {
	sql := `SELECT ur.amount_time,ur.check_time,ur.amount,ur.order_state FROM users_withdraw_records ur
	 WHERE ur.uid = ?
	ORDER BY ur.amount_time DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, uid, begin, size).QueryRows(&userDeposit)
	return
}

// 获取用户详细佣金信息总数
func GetUserDepositDetailCount(uid int) (count int, err error) {
	sql := `SELECT count(1) FROM users_withdraw_records ur
	WHERE ur.uid = ? `
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

// 根据id获取提现信息,且状态为未处理
func GetWithdrawInfoById(id int) (u *UserWithdraw, err error) {
	sql := `SELECT id,uid,order_code,amount FROM users_withdraw_records WHERE id=? AND order_state=1`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&u)
	return
}

// 获取用户的绑卡信息
func GerUsersBankcardById(Uid int) (b *UsersBankcard, err error) {
	o := orm.NewOrm()
	err = o.Raw(`SELECT id,uid,real_name,card_number,id_card,bank_mobile,create_time FROM users_bankcards WHERE uid = ?`, Uid).QueryRow(&b)
	return
}

// 将订单设置为异常订单
func SetIllegalOrder(order_code string) (err error) {
	sql := `UPDATE users_withdraw_records SET is_illegal_order=1 WHERE order_code=?`
	_, err = orm.NewOrm().Raw(sql, order_code).Exec()
	return
}

// 查询订单状态
func GetUserWithdrawDepositStatus(order_code string) (resultPay string, err error) {
	sql := `SELECT result_pay FROM users_withdraw_records WHERE order_code=?`
	err = orm.NewOrm().Raw(sql, order_code).QueryRow(&resultPay)
	return
}

// 通过订单号查询用户id和金额
func GetUidByOrderNumber(orderNo string) (uid int, amount float64, err error) {
	sql := `SELECT uid,amount FROM users_withdraw_records WHERE order_code=?`
	err = orm.NewOrm().Raw(sql, orderNo).QueryRow(&uid, &amount)
	return
}

// 修改提现订单状态
func UpdateWithdraw(result_pay, ret_code, order_code string, o orm.Ormer) (err error) {
	sql := `UPDATE users_withdraw_records SET result_pay=?,ret_code=? ,finish_time=now() WHERE order_code=?`
	_, err = o.Raw(sql, result_pay, ret_code, order_code).Exec()
	return
}

// 修改提现记录表的状态-同意
func AgreeWithdrawDeposit(orderCode string, sysId int) (err error) {
	sql := `UPDATE users_withdraw_records SET order_state=2,check_people=?,check_time=NOW() WHERE order_code=? `
	_, err = orm.NewOrm().Raw(sql, sysId, orderCode).Exec()
	return

}

// 拒绝提现
func RefuseWithdrawDeposit(userDeposit *UserWithdraw, sysId int) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
			return
		}
		o.Commit()
	}()
	//修改提现记录表的状态
	sql := `UPDATE users_withdraw_records SET order_state=3,check_people=?,result_pay=REFUSE,check_time=NOW() WHERE uid= ? AND order_code=? `
	_, err = o.Raw(sql, sysId, userDeposit.Uid, userDeposit.OrderCode).Exec()
	if err != nil {
		return
	}
	//修改钱包余额
	sql = `UPDATE users_wallet SET account_balance=account_balance+? WHERE uid=?`
	_, err = o.Raw(sql, userDeposit.Amount, userDeposit.Uid).Exec()
	if err != nil {
		return
	}
	return
}

// 校验提现次数
func CheckWithdrawCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_withdraw_records WHERE uid=? AND amount_time>CURDATE()`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

// 添加提现订单
func AddWithdrawRecord(uid int, order_code string, amount float64) (err error) {
	sql := `INSERT INTO users_withdraw_records (uid,order_code,amount,amount_time,order_state,result_pay) VALUES(?,?,?,NOW(),1,"SENDING")`
	_, err = orm.NewOrm().Raw(sql, uid, order_code, amount).Exec()
	return
}

// 获取需要主动查询的提现记录
func GetFenfaWithdrawOrder() (order []*UserWithdrawOrder, err error) {
	sql := `SELECT uid,order_code,result_pay,ret_code FROM users_withdraw_records WHERE search_time<5 AND is_illegal_order=0 AND order_state=2 AND amount_time<? AND result_pay NOT IN ("FAILURE","SUCCESS","CANCEL","REFUSE")`
	time := time.Now().Add(-time.Minute * 5)
	_, err = orm.NewOrm().Raw(sql, time).QueryRows(&order)
	return
}

// 更新订单的主动查询次数
func UpdateSearchTime(o string) (err error) {
	sql := `UPDATE users_withdraw_records SET search_time=search_time+1 WHERE order_code=?`
	_, err = orm.NewOrm().Raw(sql, o).Exec()
	return
}
