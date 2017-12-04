package services

import (
	"encoding/json"
	"fmt"
	"wr_fenfa/models"
	"wr_fenfa/utils"
)

// the service for log
func AutoInsertLogToDB() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[AutoInsertLogToDB]", err)
		}
	}()
	for {
		utils.Rc.Brpop(utils.CacheKeySystemLogs, func(b []byte) {
			var log models.SysFenfaLog
			if err := json.Unmarshal(b, &log); err != nil {
				fmt.Println("json unmarshal wrong!")
			}
			if _, err := models.AddLogs(&log); err != nil {
				fmt.Println(err.Error(), log)
			}
		})
	}
}
