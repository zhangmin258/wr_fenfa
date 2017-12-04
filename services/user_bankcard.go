package services

import (
	"wr_fenfa/models"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"wr_fenfa/utils"
)

// 聚合银行卡四元素认证
func JuheCheckBankCard(userName, bankCardNumber, idCard, account string) (result models.JuheCheckBankCardResult, err error) {
	params := make(map[string]string)
	params["key"] = "f8f3798f17148f96331d50580ffb1296" //
	params["bankcard"] = bankCardNumber                // 卡号
	params["realname"] = userName                      // 用户名
	params["idcard"] = idCard                          // 身份证
	params["mobile"] = account                         // 手机号
	url := utils.JoinGetUrl("http://v.juhe.cn/verifybankcard4/query", params)
	response, err := http.Get(url)
	if err != nil {
		return
	}
	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return
	}
	return
}
