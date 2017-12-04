package controllers

import (
	"time"
	"wr_fenfa/models"
	"wr_fenfa/cache"
	"wr_fenfa/utils"
	"strings"
	"fmt"
	"strconv"
)

// 用户数据接口
type DataController struct {
	BaseController
}

func (c *DataController) DataPage() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	user := c.User
	// 读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	condition := ""
	condition1 := ""
	params := []interface{}{}
	params1 := []interface{}{}
	startDate := c.GetString("startDate", time.Now().AddDate(0, 0, -6).Format("2006-01-02"))
	endDate:=""
	table :=""
	if user.AccountType==0 { // 外部用户
		table=" users_income_record_out "
		endDate = c.GetString("endDate", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	}else{
		table=" users_income_record "
		endDate = c.GetString("endDate", time.Now().Format("2006-01-02"))
	}

	productName := c.GetString("productName")
	userId := c.GetString("userId") // 当管理员查看指定用户的推广数据时，这个才会有值
	date := c.GetString("date")
	if date != "" {
		dates := strings.Split(date, " ")
		startDate = dates[0]
		endDate = dates[0]
	}
	if userId != "" {
		condition += " AND ui.uid=? "
		condition1 += " AND uid=? "
		params = append(params, userId)
		params1 = append(params1, userId)
	}
	if productName != "" {
		condition += " AND ui.product_name LIKE ? "
		params = append(params, "%"+productName+"%")
	}
	if startDate != "" {
		condition += " AND ui.create_date >= ? "
		params = append(params, startDate)
	}
	if endDate != "" {
		condition += " AND ui.create_date <= ? "
		params = append(params, endDate)
	}
	if user.AccountType == 0 {
		condition += " AND ui.uid=? "
		condition1 += " AND uid=? "
		params = append(params, user.Id)
		params1 = append(params1, user.Id)
	}
	// 总注册人数、七日注册人数、昨日注册人数、日均注册人数
	allRegisterCount, err := models.GetAllRegister(table,condition1, params1) // 总注册人数
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总注册人数异常", err.Error(), c.Ctx.Input)
	}
	sevenRegisterCount, err := models.GetSevenRegister(table,condition1, params1) // 七日注册人数
	if err != nil && err.Error() != utils.ErrNoRow()  {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取七日注册人数异常", err.Error(), c.Ctx.Input)
	}
	yesterdayRegisterCount, err := models.GetYesterdayRegister(table,condition1, params1) // 昨日注册人数
	if err != nil && err.Error() != utils.ErrNoRow()  {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取昨日注册人数异常", err.Error(), c.Ctx.Input)
	}
	days := 0
	count := 0
	var registerDetail []models.Registers
	if user.AccountType == 1 && userId=="" {
		days, err = models.GetAllDayAdmin() // 管理员首页总收益天数
		if err != nil{
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "管理员获取总收益天数异常", err.Error(), c.Ctx.Input)
		}
	} else if userId !=""{
		id,err:=strconv.Atoi(userId)
		if err != nil{
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id转换异常异常", err.Error(), c.Ctx.Input)
		}
		days, err = models.GetAllDayById(id) // 用户总收益天数
		if err != nil{
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总收益天数异常", err.Error(), c.Ctx.Input)
		}
	}else{
		days, err = models.GetAllDayById(user.Id) // 管理员查看用户的总收益天数
		if err != nil{
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取总收益天数异常", err.Error(), c.Ctx.Input)
		}
	}
	// 数据明细
	registerDetail, err = models.GetRegisterDetail(table,condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow()  {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取注册数据明细异常", err.Error(), c.Ctx.Input)
	}
	count, err = models.GetRegisterDetailCount(table,condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取收入数据总数异常！", err.Error(), c.Ctx.Input)
	}
	dayRegister := 0.00 // 日均收入
	if days != 0 {
		dayRegister = float64(allRegisterCount) / float64(days)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["registerDetail"] = registerDetail
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["startDate"] = startDate
	c.Data["endDate"] = endDate
	c.Data["userId"] = userId
	c.Data["dayRegister"] = fmt.Sprintf("%.2f", dayRegister)
	c.Data["allRegisterCount"] = allRegisterCount
	c.Data["sevenRegisterCount"] = sevenRegisterCount
	c.Data["yesterdayRegisterCount"] = yesterdayRegisterCount
	c.TplName = "data-statistics/promote-data.html"
}
