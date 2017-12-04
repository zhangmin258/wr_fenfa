package services

import (
	"wr_fenfa/models"
	"math/rand"
	"zcm_tools/email"
	"wr_fenfa/utils"
	"time"
)

/*
*外部用户数据
*/

func OutData() (err error) {
	// 拿出用户昨天带来的所有的收益
	allIncome, err := models.GetAllIncomeYesterday()
	if err != nil {
		email.Send("拿出用户昨天带来的所有的收益筹措", "err:"+err.Error(), utils.ToUsers, "weirong")
		return
	}
	var newIncome []models.OutUserIncome
	// 随机抽出80%的收益
	rand.Seed(time.Now().UnixNano()) // 利用当前时间的UNIX时间戳初始化rand包
	for _, v := range allIncome {
		if i := rand.Intn(100); i < 80 {
			newIncome = append(newIncome, v)
		}
	}
	// 保存80%的收益
	err = models.SaveOutIncome(newIncome)
	if err != nil {
		email.Send("保存80%的收益出错", "err:"+err.Error(), utils.ToUsers, "weirong")
		return
	}
	return
}
