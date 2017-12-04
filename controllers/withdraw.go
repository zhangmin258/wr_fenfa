package controllers

import (
	"zcm_tools/email"
	"wr_fenfa/cache"
	"wr_fenfa/utils"
	"wr_fenfa/services"
	"strconv"
	"time"
	"wr_fenfa/models"
	"encoding/json"
	"zcm_tools/uuid"
	"fmt"
)

// 提现管理接口
type WithdrawController struct {
	BaseController
}

// 审批记录
func (c *WithdrawController) ApprovalRecord() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	// 读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	// 手机号/账号
	if account := c.GetString("phone"); account != "" {
		condition += " AND u.account LIKE ? "
		params = append(params, "%"+account+"%")
	}
	// 审批过得用户信息
	condition += " AND ur.order_state !=1 "
	userDeposit, err := models.GetUserDeposit(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow() {
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取佣金列表失败", err.Error(), c.Ctx.Input)
		}
	}
	count, err := models.GetUserDepositCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDeposit"] = userDeposit
	c.TplName = "cash-management/shenpi-record.html"
}

// 待审核提现列表
func (c *WithdrawController) DepositList() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	// 读取分页信息
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	// 手机号/账号
	if account := c.GetString("phone"); account != "" {
		condition += " AND u.account LIKE ? "
		params = append(params, "%"+account+"%")
	}
	// 待审批用户信息
	condition += " AND ur.order_state =1 "
	userDeposit, err := models.GetUserDeposit(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil && err.Error() != utils.ErrNoRow() {
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取佣金列表失败", err.Error(), c.Ctx.Input)
		}
	}
	count, err := models.GetUserDepositCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "获取所有佣金数量失败", err.Error(), c.Ctx.Input)
	}
	pageCount := utils.PageCount(count, utils.PAGE_SIZE20)

	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["userDeposit"] = userDeposit
	c.TplName = "cash-management/commission-withdrawal.html"
}

// 佣金提现 通过
func (c *WithdrawController) WithdrawDeposit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	router := "withdraw/withdrawdeposit"
	sysId := c.User.Id
	jsonstr := string(c.Ctx.Input.RequestBody)
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var cid models.DepositStr
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "参数解析失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "参数解析失败"
		return
	}
	depId, err := strconv.Atoi(cid.DepId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "获取用户id失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取用户佣金id失败"
		return
	}
	// 防止重复提交
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHDRAWCASH_DEPID+strconv.Itoa(depId), 1, time.Second*5) {
		resultMap["err"] = "请不要重复提交哦"
		return
	}
	// 获取用户提现信息
	w, err := models.GetWithdrawInfoById(depId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		resultMap["err"] = "获取用户提现信息失败"
		return
	}
	if err != nil && err.Error() == utils.ErrNoRow() {
		resultMap["err"] = "没有该提现申请或申请已处理"
		return
	}
	// 获取用户银行卡信息
	bc, err := models.GerUsersBankcardById(w.Uid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		resultMap["err"] = "获取用户银行卡信息失败"
		return
	}
	if err != nil && err.Error() == utils.ErrNoRow() {
		resultMap["err"] = "该用户尚未绑卡"
		return
	}
	// 发起提现
	res, err := services.LLTrade(w.OrderCode, bc.CardNumber, bc.RealName, utils.SubFloatToString(w.Amount, 2))
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "发送请求过程发生异常", err.Error(), c.Ctx.Input)
		email.Send("发送请求过程发生异常", c.Ctx.Input.IP()+jsonstr+router+";err:"+err.Error(), utils.ToUsers, "weirong")
		resultMap["err"] = "发送请求过程发生异常!"
		return
	}
	// 发送成功
	if res.Ret_code == "0000" {
		resultMap["ret"] = 200
		resultMap["msg"] = "发送提现申请成功!"
	}
	// 作为失败处理
	if res.Ret_code != "0000" && res.Ret_code != "9999" && res.Ret_code != "4006" && res.Ret_code != "4007" && res.Ret_code != "4009" && res.Ret_code != "1002" && res.Ret_code != "2005" {
		go services.LoanResultUpdate(w.OrderCode, res.Ret_code, "", sysId)
		resultMap["ret"] = 403
		resultMap["err"] = "当前订单状态异常!已取消提现操作！"
		if res.Ret_msg != "" {
			resultMap["err"] = res.Ret_msg
		}
	}
	// 修改记录状态为成功
	err = models.AgreeWithdrawDeposit(w.OrderCode, c.User.Id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "修改记录状态为成功异常", err.Error(), c.Ctx.Input)
		email.Send("修改记录状态为成功异常", c.Ctx.Input.IP()+jsonstr+router+";err:"+err.Error(), utils.ToUsers, "weirong")
		resultMap["err"] = "修改记录状态为成功异常!"
		return
	}
	return
}

// 佣金提现 拒绝
func (c *WithdrawController) RefuseWithdrawDeposit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 403
	router := "withdraw/refusewithdrawdeposit"
	defer func() {
		c.Data["json"] = resultMap
		c.ServeJSON()
	}()
	var deposit models.DepositStr
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &deposit)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "参数解析失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "参数解析失败"
		return
	}
	depId, err := strconv.Atoi(deposit.DepId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "获取用户id失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "获取用户佣金id失败"
		return
	}
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHDRAWCASH_DEPID+strconv.Itoa(depId), 1, time.Minute) {
		resultMap["err"] = "请不要重复提交哦"
		return
	}
	//获取用户提现信息
	w, err := models.GetWithdrawInfoById(depId)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "获取用户提现信息失败", err.Error(), c.Ctx.Input)
		}
		resultMap["err"] = "获取用户提现信息失败"
		return
	}
	sysId := c.User.Id
	// 拒绝操作
	err = models.RefuseWithdrawDeposit(w, sysId)
	if err != nil {
		email.Send("拒绝提现失败", c.Ctx.Input.IP()+router+";err:"+err.Error(), utils.ToUsers, "weirong")
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "拒绝提现失败", err.Error(), c.Ctx.Input)
		resultMap["err"] = "拒绝提现失败"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "操作成功！"
	cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, router, "拒绝用户提现", "拒绝用户:"+fmt.Sprintf("%d", w.Uid)+"提现，金额为："+fmt.Sprintf("%.2f", w.Amount), c.Ctx.Input)
	return
}

// 用户发起提现，记录
// @router /withdrawcash [post]
func (this *WithdrawController) WalletWithdrawDeposit() {
	resultMap := make(map[string]interface{})
	resultMap["ret"] = 304
	//ip := this.Ctx.Input.IP()
	//router := "pay/walletwithdrawdeposit"
	defer func() {
		this.Data["json"] = resultMap
		this.ServeJSON()
	}()
	user:=this.User
	var pay models.Pay
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &pay)
	if err != nil {
		resultMap["err"] = "解析参数错误"
		return
	}
	// 判断是否可提现
	day := time.Now().Weekday().String()
	if day != utils.WITHDRAWDAY {
		resultMap["err"] = "每周只限周四提现！"
		return
	}
	pay.Uid = this.User.Id
	// 提交间隔校验
	if !utils.Rc.SetNX(utils.CACHE_KEY_ChECKWITHDRAWCASH_UID+strconv.Itoa(pay.Uid), 1, time.Second*5) {
		resultMap["err"] = "您提交的太快了，请稍后"
		return
	}
	// 当天提现次数校验
	count, err := models.CheckWithdrawCount(pay.Uid)
	if err != nil {
		resultMap["err"] = "查询用户今日提现次数错误"
		return
	}
	if count >= utils.WithdrawCount {
		resultMap["err"] = "您今日的提现次数已用完，请明天再来"
		return
	}
	// 校验用户钱包余额
	// 获取这周之前总收入
	userIncome, err := models.GetOutAllIncome(user.Id)
	if err != nil {
		cache.RecordLogs(user.Id, 0, user.Account, user.DisplayName, "withdraw/withdrawcash", "获取这周之前总收入异常", err.Error(), this.Ctx.Input)
	}
	// 获取总提现金额
	userWithdraw, err := models.GetOutUserWithdraw(user.Id)
	if err != nil {
		cache.RecordLogs(user.Id, 0, user.Account, user.DisplayName, "Withdraw/withdrawcash", "获取总提现金额异常", err.Error(), this.Ctx.Input)
	}
	// 总收入减去总提现金额
	money := userIncome.Money - userWithdraw.Money
	if pay.Money < utils.WithdrawMoney {
		resultMap["err"] = "最低提现金额为100元！"
		return
	}
	if money < pay.Money {
		resultMap["err"] = "钱包余额不足！"
		return
	}
	// 生成订单号
	orderNumber := "FFTX" + time.Now().Format("20060102") + uuid.NewUUID().Hex()
	// 提现记录表添加记录
	err = models.AddWithdrawRecord(pay.Uid, orderNumber, pay.Money)
	if err != nil {
		resultMap["err"] = "添加提现记录错误"
		return
	}
	// 提前给用户扣款
	err = models.WalletPay(pay.Uid, pay.Money)
	if err != nil {
		email.Send("用户发起提现，给用户钱包扣款出现异常", this.Ctx.Input.IP()+";err:"+err.Error(), utils.ToUsers, "weirong")
		resultMap["err"] = "用户发起提现，给用户钱包扣款出现异常！"
		return
	}
	resultMap["ret"] = 200
	resultMap["msg"] = "您的提现订单已提交，请等待人工审核"
}

// 管理员查看指定人提现记录
func (c *WithdrawController) UserWithdrawRecord() {
	xsrfToken := c.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.IsNeedTemplate()
	pageNum, _ := c.GetInt("page")
	if pageNum < 1 {
		pageNum = 1
	}
	// 获取用户信息
	userId, err := c.GetInt("userId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "id转换失败", err.Error(), c.Ctx.Input)
	}
	// 获取用户信息
	user, err := models.GetUserInfoById(userId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Account, c.User.DisplayName, "", "根据id获取用户失败", err.Error(), c.Ctx.Input)
	}
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
	c.TplName = "cash-management/admin-withdrawal-record.html"
}
