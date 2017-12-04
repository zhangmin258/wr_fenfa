package services

import (
	"github.com/astaxie/beego/toolbox"
	"wr_fenfa/utils"
)

func Task() {
	go AutoInsertLogToDB()

	// 分发平台的提现订单主动查询
	FenfaWithdrawBatchQuery := toolbox.NewTask("FenfaWithdrawBatchQuery", "0 */10 * * * * ", FenfaWithdrawBatchQuery)
	toolbox.AddTask("FenfaWithdrawBatchQuery", FenfaWithdrawBatchQuery)

	// 每天早上一点统计外部用户的收益
	outData := toolbox.NewTask("outData", "0 0 1 * * * ", OutData)
	toolbox.AddTask("outData", outData)

	// 启动定时任务
	if utils.RunMode == "release" {
		toolbox.StartTask()
	}
}
