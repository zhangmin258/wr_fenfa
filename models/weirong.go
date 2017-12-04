package models

import (
	"github.com/astaxie/beego/orm"
	"wr_fenfa/utils"
)

// 用户注册微融外链时接受的数据
type WeirongRequestForRegister struct {
	FenfaCode string // 邀请用户code
	Account   string // 被邀请手机号
}

// 用户注册微融第三方产品时接受的数据
type WeirongRequestForRegisterProduct struct {
	ProductId int    // 产品Id
	Account   string // 用户手机号
}

// 新用户在微融外链注册，user表添加一条普通用户记录
func InsertSimpleUser(newAccount string, oldUserId, level int) (int, error) {
	sql := `INSERT INTO users(parent_id,account,code,display_name,create_time,create_date,account_type,agent_level,is_used,is_agent ) VALUES(?,?,?,?,NOW(),CURDATE(),0,?,1,0)`
	result, err := orm.NewOrm().Raw(sql, oldUserId, newAccount, utils.NewUUID().Hex(), newAccount, level+1).Exec()
	if err == nil {
		id, _ := result.LastInsertId()
		return int(id), nil
	}
	return 0, err
}

// 新用户在微融外链注册，userAgent表添加一条代理关系记录
func InsertAgentRelation(pId, cId int) (err error) {
	sql := `INSERT INTO users_agent(uid,child_id,price_scale,create_time,create_date) VALUES(?,?,80,NOW(),CURDATE())`
	_, err = orm.NewOrm().Raw(sql, pId, cId).Exec()
	return
}

// 批量添加收益记录
func InsertIncomeRecord(m map[int]float64, cid, pid int, pname string) (err error) {
	sql := `INSERT INTO users_income_record (uid,child_id,product_id,product_name,income,create_date,create_time) VALUES(?,?,?,?,?,CURDATE(),NOW())`
	pre, _ := orm.NewOrm().Raw(sql).Prepare()
	defer pre.Close()
	for k, v := range m {
		pre.Exec(k, cid, pid, pname, v)
	}
	return
}

// 添加用户总收益
func AddUserIncome(m map[int]float64) (err error) {
	sql := `UPDATE users_wallet SET account_balance=account_balance+? WHERE uid=?`
	pre, _ := orm.NewOrm().Raw(sql).Prepare()
	defer pre.Close()
	for k, v := range m {
		pre.Exec(v, k)
	}
	return
}
