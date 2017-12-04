package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"wr_fenfa/models"
	"wr_fenfa/utils"
	"zcm_tools/crypt"
	"zcm_tools/http"
	"zcm_tools/pay"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

// 发送提现请求
func LLTrade(order_number, card_no, user_name, money_order string) (res *models.LLTradeResponse, err error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("api_version", "1.0")
	param.Add("sign_type", "RSA")
	param.Add("no_order", order_number)
	param.Add("dt_order", time.Now().Format("20060102150405")) //time.Now().Format("20060102150405")
	if utils.RunMode != "release" {
		money_order = "0.01"
	}
	param.Add("money_order", money_order)
	param.Add("card_no", card_no)
	param.Add("acct_name", user_name)
	param.Add("info_order", "微融")
	param.Add("flag_card", "0")
	param.Add("notify_url", utils.WR_FENFA_API_URL+utils.LL_Trade_CallBack)
	sign, err := pay.Sign(param, utils.LL_RSA_PRIVATE_KEY, pay.EncryptMD5withRsa)
	if err != nil {
		beego.Debug(err.Error())
	}
	param.Add("sign", sign)
	pay_load := crypt.LLEncrypt(param.Data(), []byte(utils.LL_PUBLIC_KEY))
	body, err := http.Post(utils.LL_Trade_API, `{"pay_load":"`+pay_load+`","oid_partner":"`+utils.LL_oid_partner+`"}`, "application/json;charset=UTF-8")
	if err != nil {
		return
	}
	json.Unmarshal(body, &res)
	if res.Ret_code == "4002" || res.Ret_code == "4003" || res.Ret_code == "4004" {
		// 发送确认消息
		res, err = LLTradeConfirm(order_number, res.Confirm_code)
	}
	return
}

// 实时付款确认
func LLTradeConfirm(order_number, confirm_code string) (res *models.LLTradeResponse, err error) {
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("api_version", "1.0")
	param.Add("sign_type", "RSA")
	param.Add("no_order", order_number)
	param.Add("confirm_code", confirm_code)
	param.Add("notify_url", utils.WR_FENFA_API_URL+utils.LL_Trade_CallBack)
	sign, err := pay.Sign(param, utils.LL_RSA_PRIVATE_KEY, pay.EncryptMD5withRsa)
	param.Add("sign", sign)
	pay_load := crypt.LLEncrypt(param.Data(), []byte(utils.LL_PUBLIC_KEY))
	body, err := http.Post(utils.LL_Trade_CONFIRM, `{"pay_load":"`+pay_load+`","oid_partner":"`+utils.LL_oid_partner+`"}`, "application/json;charset=UTF-8")
	if err != nil {
		return
	}
	json.Unmarshal(body, &res)
	fmt.Println(res.Ret_code)
	return
}

// 连连实时付款订单主动查询
func LLTradeQuery(order_number string) (*models.LLTradeQueryResponse, error) {
	fmt.Println("开始查询，订单号：", order_number)
	param := http.Values{}
	param.Add("oid_partner", utils.LL_oid_partner)
	param.Add("sign_type", "RSA")
	param.Add("api_version", "1.0")
	param.Add("no_order", order_number)
	body, err := pay.DoRequest(param, utils.LL_Trade_Query, utils.LL_RSA_PRIVATE_KEY, pay.ContentTypeJson, pay.EncryptMD5withRsa)
	if err != nil {
		return nil, err
	}
	var m *models.LLTradeQueryResponse
	err = json.Unmarshal(body, &m)
	if m == nil {
		return nil, errors.New("解析连连返回数据出错")
	}
	return m, nil
}

// 提现结果处理
func LoanResultUpdate(orderNumber, retCode, resultPay string, sysId int) (err error) {

	if utils.Rc.Lock(utils.CACHE_KEY_Loan_Accept_State_Deal+orderNumber, time.Minute) {
		defer utils.Rc.Delete(utils.CACHE_KEY_Loan_Accept_State_Deal + orderNumber)
		// 从mysql中获取提现订单状态
		orderStatus, err := models.GetUserWithdrawDepositStatus(orderNumber)
		if err != nil {
			beego.Debug(err.Error())
			return err
		}
		// 订单状态已处理
		if orderStatus == "SUCCESS" || orderStatus == "FAILURE" || orderStatus == "CANCEL" {
			return err
		}
		// 查询订单金额和用户id
		id, money, err := models.GetUidByOrderNumber(orderNumber)
		if err != nil {
			beego.Debug(err.Error())
			return err
		}
		o := orm.NewOrm()
		o.Begin()
		// 订单作失败处理
		if retCode != "0000" && retCode != "9999" && retCode != "4006" && retCode != "4007" && retCode != "4009" && retCode != "1002" && retCode != "2005" {
			resultPay = "FAILURE"
		}
		if resultPay == "FAILURE" || resultPay == "CANCEL" || resultPay == "CLOSED" || resultPay == "CHECK" {
			// 提现失败给用户钱包补款
			err = models.WalletRechargeTransaction(id, money, o)
			if err != nil {
				beego.Debug(err.Error())
				o.Rollback()
				return err
			}
		}
		// 修改提现记录
		err = models.UpdateWithdraw(resultPay, retCode, orderNumber, o)
		if err != nil {
			beego.Debug(err.Error())
			o.Rollback()
			return err
		}
		o.Commit()
	}
	return
}

// 提现订单定时查询
func FenfaWithdrawBatchQuery() (err error) {
	// 需要主动查询的充值订单
	orderList, err := models.GetFenfaWithdrawOrder()
	if err != nil {
		beego.Debug(err.Error())
		return err
	}
	if len(orderList) < 1 {
		return
	}
	fmt.Println("开始主动查询，总条数", len(orderList))
	for k, order := range orderList {
		time.Sleep(time.Second * 1)
		fmt.Println("当前第---", k+1, "---条")
		res, err := LLTradeQuery(order.OrderCode) // 发起查询
		// 更新查询次数
		err = models.UpdateSearchTime(order.OrderCode)
		if err != nil {
			beego.Debug(err.Error())
			return err
		}
		fmt.Println(res)
		// 修改订单状态
		err = LoanResultUpdate(order.OrderCode, res.Ret_code, res.Result_pay, 0)
		if err != nil {
			beego.Debug(err.Error())
			return err
		}
	}
	return
}
