package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// 分发平台用户
type User struct {
	Id          int
	ParentId    string    // 直接父代理ID
	Account     string    // 手机号（外部账户）/账户名（内部账户）
	Password    string    // 密码
	Code        string    // 用户的唯一标识
	Token       string    // 用户的登陆标识
	DisplayName string    // 管理员名称
	CreateTime  time.Time // 创建时间
	CreateDate  time.Time // 创建日期
	RoleId      int       // 角色ID
	StationId   int       // 岗位ID
	AccountType int       // 账户类型：0外部账户1内部账户
	AgentLevel  int       // 级别 * 例：1表示一级代理 2表示二级代理
	IsUsed      int       // 0：启用 1：冻结
	IsAgent     int       // 是否是代理 0 普通用户 1 代理用户
	RoleName    string    // 角色名称
}

type SysRole struct {
	Id          int
	DisplayName string
	Remark      string
	OrgId       int
}

// post请求获取权限栏
type RoleRequest struct {
	RoleId  string
	Rid     string
	Account string
	CheckId string
}

// 角色展示
type RoleShow struct {
	Id   int
	Name string
}

// 用户登陆
func Login(account, password string) (u *User, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM users WHERE  account=? and password=? `
	err = o.Raw(sql, account, password).QueryRow(&u)
	return
}

// 系统用户列表
func SysUserList(condition string, pars []interface{}, begin, count int) (list []User, err error) {
	sql := `SELECT su.*, sr.display_name role_name
			FROM users su
			LEFT JOIN role sr ON su.role_id=sr.id
			WHERE su.account_type=1 `
	sql += condition
	sql += " ORDER BY create_time DESC LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, pars, begin, count).QueryRows(&list)
	return list, nil
}

func SysUserCount(condition string, pars []interface{}) (count int, err error) {
	sql := `SELECT count(1)
			FROM users su
			LEFT JOIN role sr ON su.role_id=sr.id
			WHERE su.account_type=1 `
	sql += condition
	err = orm.NewOrm().Raw(sql, pars).QueryRow(&count)
	return
}

// 获取用户权限
func SysRoleList() (list []SysRole, err error) {
	sql := `SELECT * FROM role`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

// 根据id获取权限
func SysRoleByRid(rid int) (role *SysRole, err error) {
	sql := `SELECT * FROM role WHERE id=?`
	err = orm.NewOrm().Raw(sql, rid).QueryRow(&role)
	return
}

func (sr *SysRole) Insert(menu_ids []string) error {
	sql := `INSERT INTO role (display_name, remark, org_id)
			values(?, ?, ?)`
	o := orm.NewOrm()
	o.Begin()
	res, err := o.Raw(sql, sr.DisplayName, sr.Remark, sr.OrgId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	if len(menu_ids) > 0 {
		rid, err := res.LastInsertId()
		if err != nil {
			o.Rollback()
			return err
		}

		sql = ` INSERT INTO role_menu (role_id, menu_id) values(?, ?)`
		for i := 0; i < len(menu_ids); i++ {
			_, err = o.Raw(sql, rid, menu_ids[i]).Exec()
			if err != nil {
				break
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

func (sr *SysRole) Update(menu_ids []string) error {
	sql := `UPDATE role SET display_name=?, remark=?, org_id=?
			WHERE id=?`
	o := orm.NewOrm()
	o.Begin()
	_, err := o.Raw(sql, sr.DisplayName, sr.Remark, sr.OrgId, sr.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	sql = `DELETE FROM role_menu WHERE role_id=?`
	_, err = o.Raw(sql, sr.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	if len(menu_ids) > 0 {
		sql = ` INSERT INTO role_menu (role_id, menu_id) values(?, ?)`
		for i := 0; i < len(menu_ids); i++ {
			_, err = o.Raw(sql, sr.Id, menu_ids[i]).Exec()
			if err != nil {
				break
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

// 删除角色
func DelRole(rid int) error {
	sql := `DELETE FROM role WHERE id=?`
	o := orm.NewOrm()
	o.Begin()
	_, err := o.Raw(sql, rid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	sql = `DELETE FROM role_menu WHERE role_id=?`
	_, err = o.Raw(sql, rid).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

// 获取所有的菜单
func GetSysMenuTreeAll() (list []SysMenu, err error) {
	sql := `SELECT * FROM menu WHERE is_used=1 order by sort_index `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}

func GetSysMenuTreeByRoleId(role_id int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT DISTINCT d.* from  role_menu c INNER JOIN menu d on c.menu_id=d.id WHERE d.is_used=1 AND c.role_id =? order by sort_index "
	res := []SysMenu{}
	_, err := o.Raw(sql, role_id).QueryRows(&res)
	return res, err
}

func GetVisitMenuTree() ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT DISTINCT d.* from  role_menu c INNER JOIN menu d on c.menu_id=d.id WHERE d.is_used=1 AND d.is_private =1 order by sort_index "
	res := []SysMenu{}
	_, err := o.Raw(sql).QueryRows(&res)
	return res, err
}

// 获取所有角色列表
func GetAllRole() (roleShow []RoleShow, err error) {
	sql := `SELECT id,display_name AS name FROM role `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&roleShow)
	return
}

// 内部用户添加
func (u *User) Insert() error {
	sql := `INSERT INTO users (account, password, code, display_name,account_type,is_used, role_id, create_time,create_date)
			values(?, ?, ?, ?, ?, ?,?, NOW(),CURDATE())`
	_, err := orm.NewOrm().Raw(sql, u.Account, u.Password,u.Code, u.DisplayName,u.AccountType,u.IsUsed,u.RoleId).Exec()
	return err
}

// 内部用户修改
func (u *User) Update() error {
	sql := `UPDATE users SET  display_name=?,is_used=?, role_id=?
			WHERE id=?`
	_, err := orm.NewOrm().Raw(sql,  u.DisplayName, u.IsUsed,u.RoleId,u.Id).Exec()
	return err
}

// 删除系统用户
func DeleteUser(uid int) error {
	sql := `DELETE FROM users WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, uid).Exec()
	return err
}
