package models

import "github.com/astaxie/beego/orm"

// 获取所有用户钱包账户余额
func GetWalletBalanceAdmin() (balance float64, err error) {
	sql := `SELECT SUM(account_balance) FROM users_wallet`
	err = orm.NewOrm().Raw(sql).QueryRow(&balance)
	return
}

// 充值-事务
func WalletRechargeTransaction(uid int, money float64, o orm.Ormer) (err error) {
	sql := `UPDATE users_wallet SET account_balance=account_balance+? WHERE uid=?`
	_, err = o.Raw(sql, money, uid).Exec()
	return
}

// 用户支付，扣除余额
func WalletPay(uid int, money float64) (err error) {
	sql := `UPDATE users_wallet SET account_balance=account_balance-? WHERE uid=?`
	_, err = orm.NewOrm().Raw(sql, money, uid).Exec()
	return
}

// 初始化钱包
func InnitWallet(uid int) (err error) {
	sql := `INSERT INTO users_wallet (uid,account_balance,is_freeze) VALUES(?,0,0)`
	_, err = orm.NewOrm().Raw(sql, uid).Exec()
	return
}
