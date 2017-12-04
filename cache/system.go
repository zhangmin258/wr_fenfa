package cache

import (
	"encoding/json"
	"strconv"
	"wr_fenfa/models"
	"wr_fenfa/utils"
)

// 外部用户的菜单树
func GetOutUsersMenuTree() (m models.SysMenuList, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyOutUsersMenuTree) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyOutUsersMenuTree); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return
}

// 内部角色菜单树
func GetSysMenuTreeByRoleId(role_id int) (m models.SysMenuList, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuTree+strconv.Itoa(role_id)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyRoleMenuTree + strconv.Itoa(role_id)); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return
}

// 获取所有菜单信息
func GetCacheSysMenu() (m map[string]models.SysMenu, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySystemMenu) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeySystemMenu); err1 == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return
			}
		}
	}
	m, err = models.GetSysMenu()
	return
}

// 外部用户菜单列表
func GetOutUsersMenuList() (m []models.SysMenu, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyOutUsersMenuList) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyOutUsersMenuList); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return models.GetOutUsersMenuTree()
}

// 内部角色菜单列表
func GetSysMenuListByRoleId(role_id int) (m []models.SysMenu, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuList+strconv.Itoa(role_id)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyRoleMenuList + strconv.Itoa(role_id)); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return models.GetSysMenuTreeByStationId(role_id) // 获取外部用户菜单
}
