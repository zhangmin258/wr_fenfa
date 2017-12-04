package models

import (
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"time"
	"wr_fenfa/utils"
)

type SysMenuList []*SysMenu

func (list SysMenuList) Len() int {
	return len(list)
}

func (list SysMenuList) Less(i, j int) bool {
	return list[i].SortIndex < list[j].SortIndex
}

func (list SysMenuList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// 系统菜单
type SysMenu struct {
	Id        int         `orm:"column(id);pk"`
	TabName   string
	ImageUrl  string
	LinkUrl   string
	Remark    string
	ParentId  int
	SortIndex int
	ChildMenu SysMenuList `orm:"-"`
}

// 获取外部用户菜单
func GetOutUsersMenuTree() (res []SysMenu, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM menu WHERE is_used=1 AND is_private=1`
	_, err = o.Raw(sql).QueryRows(&res)
	return
}

// 根据岗位ID获取内部用户菜单
func GetSysMenuTreeByStationId(roleId int) (res []SysMenu, err error) {
	o := orm.NewOrm()
	sql := `SELECT m.* FROM  menu m INNER JOIN role_menu rm ON m.id=rm.menu_id WHERE m.is_used=1 AND role_id=?`
	_, err = o.Raw(sql, roleId).QueryRows(&res)
	return
}

// 外部权限修改
func MenuUpdate(menu_ids []string) error {
	// 把所有菜单置为2
	o := orm.NewOrm()
	o.Begin()
	sql := `UPDATE menu SET is_private = 2`
	_, err := o.Raw(sql).Exec()
	if err != nil {
		return err
	}
	sql = ` UPDATE menu SET is_private = 1 WHERE id = ? `
	sqlPre, err := o.Raw(sql).Prepare()
	defer func() {
		sqlPre.Close()
		if err != nil {
			o.Rollback()
			return
		} else {
			o.Commit()
		}
	}()
	if err != nil {
		return err
	}
	if len(menu_ids) > 0 {
		for _, v := range menu_ids {
			_, err = sqlPre.Exec(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取所有菜单信息
func GetSysMenu() (map[string]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT * FROM menu"
	var list []SysMenu
	_, err := o.Raw(sql).QueryRows(&list)
	m := map[string]SysMenu{}
	if err == nil && len(list) > 0 {
		for _, k := range list {
			m[k.LinkUrl] = k
		}
	}
	if data, err2 := json.Marshal(m); err == nil && err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.CacheKeySystemMenu, data, 5*time.Minute)
	}
	return m, err
}