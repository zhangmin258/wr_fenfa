package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 收益
type AgentIncome struct {
	Uid         int     //id
	AgentIncome float64 // 收益
}

// 所有收益
type Incomes struct {
	AgentIncome float64   // 收益
	CreateDate  time.Time // 日期
}

// 推广数据数据明细
type Registers struct {
	CreateDate  time.Time // 创建日期
	Account     string    // 用户手机号
	IsAgent     int       // 用户类型
	ProductName string    // 注册平台
	Income      float64   // 收益
}

// 根据用户id获取总收益
func GetAllIncomeById(table, condition string, params interface{}) (income float64, err error) {
	sql := `SELECT SUM(income) FROM `
	sql += table
	sql += ` WHERE 1=1 `
	sql += condition
	err = orm.NewOrm().Raw(sql, params).QueryRow(&income)
	return
}

// 根据用户id获取昨日收益
func GetYesterdayIncomeById(table, condition string, params interface{}) (income float64, err error) {
	sql := `SELECT SUM(income) FROM `
	sql += table
	sql += ` WHERE create_date=? `
	sql += condition
	err = orm.NewOrm().Raw(sql, time.Now().AddDate(0, 0, -1).Format("2006-01-02"), params).QueryRow(&income)
	return
}

//根据用户id获取所有收益天数
func GetAllDayById(uid int) (day int, err error) {
	sql := `SELECT TIMESTAMPDIFF(DAY,MIN(create_date),CURDATE()) FROM users_income_record_out WHERE uid = ?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&day)
	return
}

//管理员获取所有收益天数
func GetAllDayAdmin() (day int, err error) {
	sql := `SELECT TIMESTAMPDIFF(DAY,MIN(create_date),CURDATE())+1 FROM users_income_record `
	err = orm.NewOrm().Raw(sql).QueryRow(&day)
	return
}

// 按天查看收入明细
func GetIncomeEveryday(table, condition string, params []interface{}, begin, size int) (incomes []Incomes, err error) {
	sql := `SELECT SUM(income) AS agent_income,create_date FROM `
	sql += table
	sql += ` WHERE 1=1 `
	sql += condition
	sql += `GROUP BY create_date ORDER BY create_date DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&incomes)
	return
}

// 收入明细总数
func GetIncomeCount(table, condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT COUNT(1) FROM `
	sql += table
	sql += ` WHERE 1=1 `
	sql += condition
	sql += `GROUP BY create_date) a `
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 总注册人数
func GetAllRegister(table, condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT COUNT(1) FROM `
	sql += table
	sql += ` ui
	LEFT JOIN users u
	ON ui.uid = u.id
	WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY child_id,product_id) a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 七日注册人数
func GetSevenRegister(table, condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT COUNT(1) FROM `
	sql += table
	sql += ` ui
	LEFT JOIN users u
	ON ui.uid = u.id
	WHERE ui.create_date>=? AND ui.create_date <=? `
	sql += condition
	sql += ` GROUP BY child_id,product_id) a`
	err = orm.NewOrm().Raw(sql, time.Now().AddDate(0, 0, -6).Format("2006-01-02"), time.Now(), params).QueryRow(&count)
	return
}

// 昨日注册人数
func GetYesterdayRegister(table, condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT COUNT(1) FROM `
	sql += table
	sql += ` ui
	LEFT JOIN users u
	ON ui.uid = u.id
	WHERE ui.create_date =? `
	sql += condition
	sql += ` GROUP BY child_id,product_id) a`
	err = orm.NewOrm().Raw(sql, time.Now().AddDate(0, 0, -1).Format("2006-01-02"), params).QueryRow(&count)
	return
}

// 根据用户按天获取注册数据明细
func GetRegisterDetail(table, condition string, params []interface{}, begin, size int) (registers []Registers, err error) {
	sql := `SELECT ui.create_date,u.is_agent,u.account_type,ui.product_name,SUM(ui.income) AS income,u.account FROM `
	sql += table
	sql += ` ui
	LEFT JOIN users u
	ON ui.child_id = u.id WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY ui.child_id,product_id ORDER BY ui.create_date DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&registers)
	return
}

// 根据用户按天获取注册数据明细
func GetRegisterDetailCount(table, condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT ui.child_id FROM `
	sql += table
	sql += ` ui
	LEFT JOIN users u
	ON ui.child_id = u.id
	WHERE 1=1 `
	sql += condition
	sql += ` GROUP BY ui.child_id,product_id) a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}
