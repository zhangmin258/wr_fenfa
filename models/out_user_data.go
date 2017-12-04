package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// 用户收益
type OutUserIncome struct {
	Uid         int       // 用户id-父id
	ChildId     int       // 子id
	ProductId   int       // 产品id
	CreateDate  time.Time // 创建日期
	CreateTime  time.Time // 创建时间
	ProductName string    // 注册平台
	Income      float64   // 收益
}

// 用户余额
type UserMoney struct {
	Uid       int     // 用户id
	Money     float64 // 余额
}

// 拿出用户昨天带来的所有的收益
func GetAllIncomeYesterday() (income []OutUserIncome, err error) {
	sql := `SELECT uid,child_id,product_id,product_name,income,create_date,create_time FROM users_income_record WHERE create_date = ? `
	_, err = orm.NewOrm().Raw(sql, time.Now().AddDate(0, 0, -1).Format("2006-01-02")).QueryRows(&income)
	return
}

// 保存外部用户收益
func SaveOutIncome(newIncome []OutUserIncome) error {
	sql := `INSERT INTO users_income_record_out (uid,child_id,product_id,product_name,income,create_date,create_time)VALUES(?,?,?,?,?,?,?)`
	insertPre, err := orm.NewOrm().Raw(sql).Prepare()
	if err != nil {
		return err
	}
	defer insertPre.Close()
	for _, v := range newIncome {
		_, err = insertPre.Exec(v.Uid, v.ChildId, v.ProductId, v.ProductName, v.Income, v.CreateDate, v.CreateTime)
		if err != nil {
			return err
		}
	}
	return err
}

// 外部获取总收入
func GetOutAllIncome(id int ) (income UserMoney, err error) {
	sql := `SELECT uid,SUM(income) AS money FROM users_income_record_out WHERE YEARWEEK(date_format(create_date,'%Y-%m-%d'),1) < YEARWEEK(now(),1) AND uid = ? `
	err = orm.NewOrm().Raw(sql,id).QueryRow(&income)
	return
}

// 外部获取总提现金额
func GetOutUserWithdraw(id int) (userMoney UserMoney, err error) {
	sql := `SELECT uid,SUM(amount) AS money FROM users_withdraw_records WHERE order_state!=3 AND uid = ? `
	err=orm.NewOrm().Raw(sql,id).QueryRow(&userMoney)
	return
}