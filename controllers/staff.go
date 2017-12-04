package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/cache"
	"wr_fenfa/utils"
	"encoding/json"
	"strconv"
)

// 员工接口模块
type StaffController struct {
	BaseController
}

func (c *StaffController) StaffList() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	// 读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	user := c.User
	condition := ""
	params := []interface{}{}
	if user.AccountType == 0 { // 外部用户
		condition += ` AND (u.parent_id = ? OR u.id=? ) `
		params = append(params, user.Id)
		params = append(params, user.Id)
	}
	account := c.GetString("account")
	// 读取搜索信息
	if account != "" {
		condition += ` AND u.account LIKE ? `
		params = append(params, "%"+account+"%")
	}
	// 获取员工
	staff, err := models.GetStaff(user.AccountType,condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取员工信息失败", err.Error(), c.Ctx.Input)
		}
	count, err := models.GetStaffCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取员工总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["staff"] = staff
	c.Data["user"] = user
	c.Data["pageSize"] = utils.PAGE_SIZE20
	c.TplName = "staff-management/staff-list.html"
}

// 禁用和启用登录
func (c *StaffController) LockLogin() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	// 禁用登录的员工id
	var staffRemark models.StaffRemark
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &staffRemark)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "禁用和启用登录解析参数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析参数失败"
		return
	}
	if staffRemark.Code == "" {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "员工获取失败", "", c.Ctx.Input)
		resultMap["err"] = "员工获取失败"
		return
	}
	// 设置员工状态为禁用
	err = models.UpdateUserState(staffRemark.Code, staffRemark.State)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "设置员工状态出错", err.Error(), c.Ctx.Input)
		resultMap["err"] = "设置员工状态出错"
		return
	}
	cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "设置员工状态", "用户code："+c.User.Code+"，设置用户code："+staffRemark.Code+"的状态码为："+strconv.Itoa(staffRemark.State), c.Ctx.Input)
	resultMap["ret"] = 200
}

// 编辑员工的姓名和分成
func (c *StaffController) SaveStaffInfo() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var staffRemark models.StaffRemark
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &staffRemark)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "编辑员工的姓名和分成解析参数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析参数失败"
		return
	}
	if staffRemark.Code == "" {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "员工code不正确", "", c.Ctx.Input)
		resultMap["err"] = "员工code不正确"
		return
	}
	user := c.User
	u, err := models.GetUserInfoByCode(staffRemark.Code)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取员工信息失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取员工信息失败"
		return
	}
	staffRemark.ChildId = u.Id
	// 员工分成
	priceScale := 0
	if staffRemark.PriceScale != "" {
		priceScale, err = strconv.Atoi(staffRemark.PriceScale)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取员工分成失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "获取员工分成失败"
			return
		}
	}
	if staffRemark.UserAgentId == "" {
		// 新增员工的备注和分成
		err = models.InsertStaffRemark(user.Id, staffRemark.ChildId, priceScale, staffRemark.Remark)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "新增员工的备注和分成失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "新增员工的备注和分成失败"
			return
		}
	} else {
		// 修改员工的备注和分成
		err = models.UpdateStaffRemark(user.Id, staffRemark.ChildId, priceScale, staffRemark.Remark)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "修改员工的备注和分成失败", err.Error(), c.Ctx.Input)
			resultMap["err"] = "修改员工的备注和分成失败"
			return
		}
	}
	resultMap["ret"] = 200
}
