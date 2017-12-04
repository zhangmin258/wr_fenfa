package models

import (
	"github.com/astaxie/beego/orm"
	"wr_fenfa/utils"
	"encoding/json"
	"time"
)

// 用户验证信息
type UserRequest struct {
	User
	Vcode string // 手机验证码
	Type  int    // 用于检测更改密码操作类型，1忘记密码，2，修改密码
}

// 用户个人信息
type UserInfo struct {
	Account     string // 手机号
	Code        string // 用户唯一标识
	Url         string // 专属链接
	Name        string // 姓名
	BankAccount string // 银行预留手机号
	IdCard      string // 身份证号码
	BankCard    string // 银行卡
}

// 用户验证信息
type UserIdAndAgentLevel struct {
	Uid        int
	AgentLevel int
}

type UserPrice struct {
	ParentId   int // 上级代理id
	PriceScale int // 该子代理的分成系数
}

type AgentUserInfo struct {
	Id         int
	Account    int // 上级代理code
	IsAgent    int // 是否是代理 0 普通用户 1 代理用户
	PriceScale int // 代理系数
	AgentLevel int // 级别 * 例：1表示一级代理 2表示二级代理
}

// 新增用户
func AddUsers(m *User) (id int, err error) {
	o := orm.NewOrm()
	sql := `INSERT INTO users (code,account,password,display_name,create_time,create_date,account_type,agent_level,is_used,is_agent,agent_time) VALUES(?,?,?,?,now(),CURRENT_DATE(),?,?,?,?,now())`
	result, err := o.Raw(sql, m.Code, m.Account, m.Password, m.DisplayName, m.AccountType, m.AgentLevel, m.IsUsed, m.IsAgent).Exec()
	id2, err := result.LastInsertId()
	if err != nil {
		return
	}
	id = int(id2)
	return
}

// 将普通用户升级成为代理用户
func UpdateSimpleUsersBeAgent(m *User) error {
	o := orm.NewOrm()
	sql := `UPDATE users SET password=?,is_used=0,is_agent=1 ,agent_time=NOW() WHERE account=?`
	_, err := o.Raw(sql, m.Password, m.Account).Exec()
	return err
}

// 检测是否存在该用户
func CheckUser(account string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE account = ?`
	err = orm.NewOrm().Raw(sql, account).QueryRow(&count)
	return
}

// 检测是否存在该代理用户
func CheckAgentUser(account string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE account = ? AND is_agent=1`
	err = orm.NewOrm().Raw(sql, account).QueryRow(&count)
	return
}

// 修改用户信息
func UpdateUser(user User) error {
	sql := `UPDATE users SET password =? WHERE account = ? `
	_, err := orm.NewOrm().Raw(sql, utils.MD5(user.Code+user.Password), user.Account).Exec()
	return err
}

// 根据用户code查询用户的id和代理级别
func GetUsersIdAndLevelByCode(code string) (r *UserIdAndAgentLevel, err error) {
	sql := `SELECT id as uid,agent_level FROM users WHERE code=? `
	err = orm.NewOrm().Raw(sql, code).QueryRow(&r)
	return
}

// 根据用户account查询用户code
func GetCodeByAccount(account string) (code string, err error) {
	sql := `SELECT code  FROM users WHERE account=? `
	err = orm.NewOrm().Raw(sql, account).QueryRow(&code)
	return
}

// 根据code获取用户详细信息
func GetUserInfoByCode(code string) (u *User, err error) {
	sql := `SELECT * FROM users WHERE code=? `
	err = orm.NewOrm().Raw(sql, code).QueryRow(&u)
	if err == nil {
		if data, err2 := json.Marshal(u); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKeyUserInfo+code, data, 5*time.Minute)
		}
	}
	return
}

// 根据id获取用户详细信息
func GetUserInfoById(id int) (u *User, err error) {
	sql := `SELECT * FROM users WHERE id=? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&u)
	return
}

// 根据account获取用户详细信息
func GetUserInfoByAccount(account string) (u *User, err error) {
	sql := `SELECT * FROM users WHERE account=? `
	err = orm.NewOrm().Raw(sql, account).QueryRow(&u)
	return
}

// 查看是否有该用户
func CheckIsRegiste(account string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE account=? `
	err = orm.NewOrm().Raw(sql, account).QueryRow(&count)
	return
}

// 根据uid获取用户详细信息
func GetUserInfoByUid(uid int) (u *User, err error) {
	sql := `SELECT * FROM users WHERE id=? `
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&u)
	return
}

// 获取用户详细信息
func GetUserInfo(id int) (user UserInfo, err error) {
	sql := `SELECT u.account,u.code,ub.card_number AS bank_card,ub.id_card,ub.bank_mobile AS bank_account,ub.real_name AS name FROM users u
	LEFT JOIN users_bankcards ub
	ON u.id = ub.uid WHERE u.id =? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&user)
	return
}

func UpdateUserState(code string, state int) error {
	sql := `UPDATE users SET is_used = ? WHERE code = ? `
	_, err := orm.NewOrm().Raw(sql, state, code).Exec()
	return err
}

// 根据用户的code 查出 用户的分成系数和上级代理
func QueryParentAndPrice(id int) (userPrice UserPrice, err error) {
	// 获取上级代理给的分成系数
	sql := `SELECT uid AS parent_id,price_scale FROM users_agent WHERE child_id= ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&userPrice)
	return
}

// 根据用户的code 查出 用户的分成系数和上级代理
func GetParent(id int) (pid int, err error) {
	// 获取上级代理给的分成系数
	sql := `SELECT uid AS parent_id FROM users_agent WHERE child_id= ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&pid)
	return
}

// 根据用户的code 查出 用户的分成系数和上级代理
func QueryParentscale(id int) (priceScale int, err error) {
	// 获取上级代理给的分成系数
	sql := `SELECT price_scale FROM users_agent WHERE child_id= ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&priceScale)
	return
}

// 根据account查询用户uid
func GetUidByAccount(account string) (uid int, err error) {
	sql := `SELECT id AS uid FROM users WHERE account=?`
	err = orm.NewOrm().Raw(sql, account).QueryRow(&uid)
	return
}

// 查询用户是否是代理
func SearchIsAgent(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE id=? AND is_agent=1`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

// 查询用户是否是分发用户
func IsFenfaUser(account string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE account=?`
	err = orm.NewOrm().Raw(sql, account).QueryRow(&count)
	return
}
