package models

import (
	"wr_fenfa/utils"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDataBase("default", "mysql", utils.MYSQL_FENFA_URL)
	orm.RegisterDataBase("wr_log", "mysql", utils.MYSQL_LOG_URL)
	orm.RegisterDataBase("wr_backup", "mysql", utils.MYSQL_BACKUP_URL)
	orm.RegisterDataBase("wr", "mysql", utils.MYSQL_URL)
	orm.RegisterModel(
		new(SysFenfaLog))
}
