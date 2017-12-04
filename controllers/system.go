package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/cache"
	"wr_fenfa/utils"
	"strings"
	"wr_fenfa/services"
	"strconv"
	"encoding/json"
)

// 内部用户及权限框架接口
type SystemController struct {
	BaseController
}

func (c *SystemController) RoleList() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	rolelist, err := models.SysRoleList()
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取系统角色失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = rolelist
	c.TplName = "system-management/role_list.html"
}

func (c *SystemController) RoleEdit() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	rid_str := c.GetString("rid")
	var role *models.SysRole
	if rid_str == "" {
		role = &models.SysRole{}
	} else {
		rid, err := c.GetInt("rid")
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id错误！", err.Error(), c.Ctx.Input)
			c.Ctx.WriteString("id错误")
			return
		}
		role, err = models.SysRoleByRid(rid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "根据id获取系统角色失败！", err.Error(), c.Ctx.Input)
			}
			c.Ctx.WriteString(err.Error())
			return
		}
	}
	c.Data["role"] = role
	c.TplName = "system-management/role_edit.html"
}

func (c *SystemController) VisitRoleEdit() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	c.TplName = "system-management/visit_role_edit.html"
}

func (c *SystemController) RoleAdd() {
	defer c.ServeJSON()
	var role *models.SysRole
	var err error
	var roleRequest models.RoleRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &roleRequest)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "解析参数失败"}
		return
	}
	rid, err := strconv.Atoi(roleRequest.Rid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id转换错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id转换错误"}
		return
	}
	if rid == 0 {
		role = &models.SysRole{}
	} else {
		role, err = models.SysRoleByRid(rid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "根据id获取系统角色失败！", err.Error(), c.Ctx.Input)
			}
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
			return
		}
	}
	ids_str := roleRequest.CheckId // 菜单权限
	ids := strings.Split(ids_str, ",")
	role.DisplayName = roleRequest.Account
	if rid == 0 {
		err = role.Insert(ids)
	} else {
		err = role.Update(ids)
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "新增或修改角色失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	//删除缓存
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuTree+strconv.Itoa(rid)) {
		utils.Rc.Delete(utils.CacheKeyRoleMenuTree + strconv.Itoa(rid))
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//删除角色
func (c *SystemController) DelRole() {
	defer c.ServeJSON()
	var roleRequest models.RoleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleRequest)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "解析参数失败"}
		return
	}
	rid, err := strconv.Atoi(roleRequest.Rid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	if rid == 1 { // 超级管理员不给删
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	err = models.DelRole(rid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "删除角色失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *SystemController) MenuData() {
	defer c.ServeJSON()
	var roleRequest models.RoleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleRequest)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "解析参数失败"}
		return
	}
	var list []models.SysMenu
	if roleRequest.RoleId == "all" { // 所有菜单
		list, err = models.GetSysMenuTreeAll()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取所有菜单错误", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "获取所有菜单错误"}
			return
		}
	} else if roleRequest.RoleId == "visit" { // 外部菜单
		list, err = models.GetVisitMenuTree()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取外部用户有的菜单错误", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "获取外部用户有的菜单错误"}
			return
		}
	} else {
		// 该角色有的菜单
		rid, err := strconv.Atoi(roleRequest.Rid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id获取失败", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id获取失败"}
			return
		}
		list, err = models.GetSysMenuTreeByRoleId(rid)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取该角色有的菜单错误", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "获取该角色有的菜单错误"}
			return
		}
	}
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取目录失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	m := services.GetSysMenuZTree(list)
	c.Data["json"] = map[string]interface{}{"ret": 200, "m": m}
}

//外部用户权限编辑
func (c *SystemController) VisitMenuEdit() {
	defer c.ServeJSON()
	var roleRequest models.RoleRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &roleRequest)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "解析参数失败"}
		return
	}
	ids_str := roleRequest.CheckId // 菜单权限
	ids := strings.Split(ids_str, ",")
	err = models.MenuUpdate(ids)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "修改外部用户菜单错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "修改外部用户菜单错误"}
		return
	}
	//删除缓存
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyOutUsersMenuTree) {
		utils.Rc.Delete(utils.CacheKeyOutUsersMenuTree)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}
