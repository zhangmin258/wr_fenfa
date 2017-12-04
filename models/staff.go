package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// 员工信息
type Staff struct {
	Code        string    // 用户唯一标识
	DisplayName string    // 用户名称
	CreateTime  time.Time // 创建时间
	Account     string    // 手机号
	IsUsed      int       // 是否在使用
	PriceScale  int       // 分成系数
	UserAgentId int       // 代理id
	AgentLevel  int       // 代理级别
}

// 员工搜索信息
type StaffSeacher struct {
	Page    int    // 页码
	Account string // 手机号
}

// 员工的姓名和分成
type StaffRemark struct {
	UserAgentId string // 代理id
	ChildId     int    // 子代理id
	Code        string
	Remark      string
	PriceScale  string
	State       int //0 启用 1 禁用
}

// 获取员工
func GetStaff(accountType int, condition string, params []interface{}, begin, size int) (staff []Staff, err error) {
	sql := `SELECT code,u.display_name,u.create_time,u.account,is_used,price_scale,u.agent_level FROM users u
			LEFT JOIN users_agent a
			ON u.id = a.child_id WHERE is_agent=1 `
	if accountType == 0 {
		sql = `SELECT a.id AS user_agent_id,code,a.remark AS display_name,u.agent_time AS create_time,u.account,is_used,price_scale,u.agent_level
		FROM users u LEFT JOIN users_agent a
		ON u.id = a.child_id
		WHERE u.is_agent=1 `
	}
	sql += condition
	sql += `ORDER BY u.agent_time DESC LIMIT ?, ? `
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&staff)
	return
}

/*// 内部用户根据id获取员工
func GetAdminStaffById(condition string, params []interface{}, begin, size int) (staff []Staff, err error) {
	sql := `SELECT code,u.display_name,u.create_time,u.account,is_used,price_scale,u.agent_level FROM users u
	LEFT JOIN users_agent ua
	ON u.id = ua.child_id WHERE is_agent=1 `
	sql += condition
	sql += `ORDER BY u.create_date DESC LIMIT ?, ? `
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&staff)
	return
}*/

// 根据id获取员工
func GetStaffCount(condition string, params []interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u
	LEFT JOIN users_agent a
	ON u.id = a.child_id WHERE is_agent=1 `
	sql += condition
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 修改员工的备注及分成
func UpdateStaffRemark(pid, childId, priceScale int, remark string) error {
	sql := `UPDATE users_agent SET price_scale = ?,remark = ? WHERE uid = ? AND child_id = ? `
	_, err := orm.NewOrm().Raw(sql, priceScale, remark, pid, childId).Exec()
	return err
}

// 新增员工的备注和分成
func InsertStaffRemark(pid, childId, priceScale int, remark string) error {
	sql := `INSERT INTO users_agent (uid,child_id,price_scale,remark,create_time,create_date) VALUES (?,?,?,?,NOW(),CURDATE())`
	_, err := orm.NewOrm().Raw(sql, pid, childId, priceScale, remark).Exec()
	return err
}
