package controllers

import (
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"wr_fenfa/cache"
	"fmt"
	"time"
)

// 佣金结算接口
type CommissionController struct {
	BaseController
}

// 价格列表
func (c *CommissionController) ProductPriceList() {
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
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "从微融获取上线的产品总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	user := c.User
	// 获取产品的分成系数
	if user.AccountType == 0 {
		// 获取用户的代理级别
		r, err := models.GetUsersIdAndLevelByCode(user.Code)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取用户的代理级别异常！", err.Error(), c.Ctx.Input)
		}
		if r.AgentLevel > 1 {
			price := 0.00
			for i := 0; i < r.AgentLevel-1; i++ {
				// 根据用户的id 查出 用户的分成系数和上级代理
				userPrice, err := models.QueryParentAndPrice(r.Uid)
				if err != nil {
					cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取用户的分成系数和上级代理异常！", err.Error(), c.Ctx.Input)
				}
				if i == 0 {
					price = float64(userPrice.PriceScale) / 100
				} else {
					price = price * float64(userPrice.PriceScale) / 100
				}
				r.Uid = userPrice.ParentId
			}
			for k, v := range products {
				products[k].AgentPrice = v.AgentPrice * price
			}
		}
	}
	c.Data["products"] = products
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["pageNum"] = pageNum
	c.Data["user"] = user
	c.TplName = "commission-settlement/price-list.html"
}

// 佣金明细
func (c *CommissionController) CommissionInfo() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	// 读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	condition := ""
	params := []interface{}{}
	condition1 := ""
	params1 := []interface{}{}
	user := c.User
	startDate := c.GetString("startDate", time.Now().AddDate(0, 0, -6).Format("2006-01-02"))
	endDate := ""
	table := ""
	if user.AccountType == 0 { // 外部用户
		table = " users_income_record_out "
		endDate = c.GetString("endDate", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	} else {
		table = " users_income_record "
		endDate = c.GetString("endDate", time.Now().Format("2006-01-02"))
	}
	if user.AccountType == 0 {
		condition1 = " AND uid= ? "
		condition += " AND uid= ? "
		params = append(params, user.Id)
		params1 = append(params1, user.Id)
	}
	if startDate != "" {
		condition += " AND create_date >= ?"
		params = append(params, startDate)
	}
	if endDate != "" {
		condition += " AND create_date <= ? "
		params = append(params, endDate)
	}
	// 总收入、昨日收入、日均收入、可提现收入
	allIncome, err := models.GetAllIncomeById(table, condition1, params1) // 总收入
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总收入异常", err.Error(), c.Ctx.Input)
	}
	yesterdayIncome, err := models.GetYesterdayIncomeById(table, condition1, params1) // 昨日收入
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取昨日收入异常", err.Error(), c.Ctx.Input)
	}
	days := 0
	// 用户的银行信息
	var bankCard models.Card
	var income float64
	if user.AccountType == 0 {
		days, err = models.GetAllDayById(user.Id) // 总收益天数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总收益天数异常", err.Error(), c.Ctx.Input)
		}
		bankCard, err = models.GetBankCardById(user.Id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取用户银行信息异常", err.Error(), c.Ctx.Input)
		}
		// 可提现收入
		// 获取这周之前总收入
		userIncome, err := models.GetOutAllIncome(user.Id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取这周之前总收入异常", err.Error(), c.Ctx.Input)
		}
		// 获取总提现金额
		userWithdraw, err := models.GetOutUserWithdraw(user.Id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总提现金额异常", err.Error(), c.Ctx.Input)
		}
		// 总收入减去总提现金额
		income = userIncome.Money - userWithdraw.Money
	} else {
		days, err = models.GetAllDayAdmin() // 管理员总收益天数
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "管理员获取总收益天数异常", err.Error(), c.Ctx.Input)
		}
		income, err = models.GetWalletBalanceAdmin() // 管理员查看所有可提现收入
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取可提现收入异常", err.Error(), c.Ctx.Input)
		}
	}
	dayIncome := 0.00 // 日均收入
	if days != 0 {
		dayIncome = allIncome / float64(days)
	}
	// 数据明细
	incomeData, err := models.GetIncomeEveryday(table, condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取收入数据明细异常", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetIncomeCount(table, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取收入数据总数异常！", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["income"] = fmt.Sprintf("%.2f", income)
	c.Data["allIncome"] = fmt.Sprintf("%.2f", allIncome)
	c.Data["dayIncome"] = fmt.Sprintf("%.2f", dayIncome)
	c.Data["yesterdayIncome"] = fmt.Sprintf("%.2f", yesterdayIncome)
	c.Data["startDate"] = startDate
	c.Data["endDate"] = endDate
	c.Data["incomeData"] = incomeData
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["accountType"] = user.AccountType
	c.Data["bankCard"] = bankCard
	c.TplName = "commission-settlement/commission-detail.html"
}

// 提现记录
func (c *CommissionController) WithdrawRecord() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	// 获取用户信息
	user := c.User
	userId := user.Id
	var err error
	userDepositDetail, err := models.GetUserDepositDetail(userId, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取用户提现明细失败", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetUserDepositDetailCount(userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["user"] = user
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDepositDetail"] = userDepositDetail
	c.TplName = "commission-settlement/withdrawal-record.html"
}
