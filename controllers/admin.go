package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"wr_fenfa/cache"
	"encoding/json"
	"strconv"
)

// 管理员接口
type AdminController struct {
	BaseController
}

// 内部用户列表
func (c *AdminController) UserList() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}
	condition := ""
	pars := []interface{}{}
	if account := c.GetString("account"); account != "" {
		condition += " AND su.account=? "
		pars = append(pars, account)
	}
	if username := c.GetString("user"); username != "" {
		condition += " AND su.display_name=? "
		pars = append(pars, username)
	}
	roleId, err := c.GetInt("post",0)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取角色id失败", err.Error(), c.Ctx.Input)
	}
	if roleId != 0 {
		condition += " AND sr.id=? "
		pars = append(pars, roleId)
	}
	accountState, err := c.GetInt("account-state",0)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取账号状态失败", err.Error(), c.Ctx.Input)
	}
	if accountState != 0 {
		condition += " AND su.is_used=? "
		pars = append(pars, accountState-1)
	}

	list, err := models.SysUserList(condition, pars, utils.StartIndex(page, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error()!=utils.ErrNoRow(){
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取系统用户列表失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.SysUserCount(condition, pars)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取系统用户总数失败", err.Error(), c.Ctx.Input)
	}
	pagecount := utils.PageCount(count, utils.PAGE_SIZE20)
	// 获取所有角色列表
	role, err := models.GetAllRole()
	if err != nil && err.Error()!=utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取所有角色列表失败", err.Error(), c.Ctx.Input)
	}
	c.Data["role"] = role
	c.Data["currpage"] = page
	c.Data["pagecount"] = pagecount
	c.Data["list"] = list
	c.Data["count"] = count
	c.TplName = "system-management/user-system.html"
}

// 产品列表
func (c *AdminController) PriceList() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	condition := ``
	params := []string{}
	// 读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	// 产品名称
	name := c.GetString("name")
	if name != "" {
		condition += ` AND name LIKE ?`
		params = append(params, "%"+name+"%")
	}
	// 从微融获取上线的产品
	products, err := models.GetProductList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "从微融获取上线的产品异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetProductCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取员工总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["products"] = products
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["pageNum"] = pageNum
	c.TplName = "system-management/price-set.html"
}

// 保存产品价格
func (c *AdminController) SaveProduct() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var productEdit models.ProductEdit
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &productEdit)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "解析参数失败"
		return
	}
	err = models.UpdateProduct(productEdit)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "修改产品价格失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "修改产品价格失败"
		return
	}
	cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "修改产品价格", "产品id："+strconv.Itoa(productEdit.ProductId)+"修改价格为："+strconv.FormatFloat(productEdit.AgentPrice, 'f', -1, 64), c.Ctx.Input)
	resultMap["ret"] = 200
}

// 编辑系统用户
func (c *AdminController) UserEdit() {
	defer c.ServeJSON()
	var user *models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "参数解析失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "参数解析失败", "err": err.Error()}
	}
	// 初始化用户
	user.AccountType = 1
	user.Code = utils.NewUUID().Hex()
	user.Password = utils.MD5(user.Code + user.Password)
	if user.Id == 0 {
		err = user.Insert()
	} else {
		err = user.Update()
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "添加系统用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *AdminController) DelUser() {
	defer c.ServeJSON()
	var user *models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "参数解析失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": "参数解析失败", "err": err.Error()}
	}
	if user.Id == 1 { // 系统管理员不给删
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	err = models.DeleteUser(user.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "删除系统用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}
