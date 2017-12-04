package services

import (
	"encoding/json"
	"sort"
	"wr_fenfa/cache"
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"strconv"
)

//获取用户菜单
func GetSysMenuTree(user *models.User) (models.SysMenuList, error) {
	var menu models.SysMenuList
	var m []models.SysMenu
	if user.AccountType == 0 { // 外部用户,获取所有非内部tab
		t, _ := cache.GetOutUsersMenuTree()
		if t != nil {
			return t, nil
		}
		m, _ = models.GetOutUsersMenuTree() // 缓存失效数据库查询
		menu = MenuFactor(m, menu)
		if data, err2 := json.Marshal(menu); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKeyOutUsersMenuTree, data, utils.RedisCacheTime_Role)
		}
	} else { // 内部用户，根据RoleId获取所有tab
		t, _ := cache.GetSysMenuTreeByRoleId(user.RoleId)
		if t != nil {
			return t, nil
		}
		m, _ := models.GetSysMenuTreeByStationId(user.RoleId) //根据岗位ID,获取菜单信息
		menu = MenuFactor(m, menu)
		if data, err2 := json.Marshal(menu); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKeyRoleMenuTree+strconv.Itoa(user.RoleId), data, utils.RedisCacheTime_Role)
		}
	}
	return menu, nil
}

// 加工menu
func MenuFactor(m []models.SysMenu, menu models.SysMenuList) (models.SysMenuList) {
	l := len(m)
	for i := 0; i < l; i++ {
		if m[i].ParentId == 0 {
			for j := 0; j < l; j++ {
				if m[j].ParentId == m[i].Id {
					m[i].ChildMenu = append(m[i].ChildMenu, &m[j])
				}
			}
			sort.Sort(m[i].ChildMenu)
			menu = append(menu, &m[i])
		}
	}
	sort.Sort(menu)
	return menu
}

func GetSysMenuZTree(list []models.SysMenu) []map[string]interface{} {
	var menu []map[string]interface{}
	l := len(list)
	if l == 0 {
		return []map[string]interface{}{}
	}
	for i := 0; i < l; i++ {
		menu = append(menu, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].TabName})
	}
	return menu
}

// 校验菜单权限
func CheckMenu(url string, u *models.User) (b bool, err error) {
	b = false
	allMenus, err := cache.GetCacheSysMenu() // 获取所有菜单列表
	if err != nil {
		return
	}
	if _, ok := allMenus[url]; ok {
		var menuList []models.SysMenu // 获取当前用户可用菜单列表
		if u.AccountType == 0 {// 外部用户
			menuList, err = cache.GetOutUsersMenuList()
			if err != nil {
				return
			}
		} else { // 内部用户
			menuList, err = cache.GetSysMenuListByRoleId(u.RoleId)
			if err != nil {
				return
			}
		}
		for _, v := range menuList { // 判断用户是否有该地址的权限
			if v.LinkUrl == url {
				b = true
			}
		}
	}else{
		b = true
	}
	return
}
