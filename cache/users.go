package cache

import (
	"encoding/json"
	"wr_fenfa/utils"
	"wr_fenfa/models"
)

//===========================获取用户的基本信息缓存====================================
// 通过uid或者account获取用户的【注册信息】缓存
/*func GetUsersByAccountCache(account string) (m *models.User, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.KEY_USERS_IFNO_REGISTER+"_"+account) {
		if data, err1 := utils.Rc.RedisBytes(utils.KEY_USERS_IFNO_REGISTER + "_" + account); err1 == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return
			}
		}
	}
	return models.GetUserInfoByAccount(account)
}
func GetUsersByIdCache(uid int) (m *models.User, err error) {
	if utils.Re == nil && uid != 0 && utils.Rc.IsExist(utils.KEY_USERS_IFNO_REGISTER+"_"+strconv.Itoa(uid)) {
		if data, err1 := utils.Rc.RedisBytes(utils.KEY_USERS_IFNO_REGISTER + "_" + strconv.Itoa(uid)); err1 == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return
			}
		}
	}
	return models.GetUserInfoByUid(uid)
}*/

//  通过code获取用户的【注册信息】缓存
func GetUserByCodeCache(code string) (m *models.User, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyUserInfo+code) {
		if date, err1 := utils.Rc.RedisBytes(utils.CacheKeyUserInfo + code); err1 == nil {
			err = json.Unmarshal(date, &m)
			if m != nil {
				return
			}
		}
	}
	return models.GetUserInfoByCode(code)
}
