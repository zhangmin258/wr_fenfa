package services

import (
	"wr_fenfa/models"
)

// 向用户表里添加完成外链注册的用户（普通用户）
func InsertSimpleUser(oldCode, newAccount string) (err error) {
	// 查询邀请人ID和代理级别
	r, err := models.GetUsersIdAndLevelByCode(oldCode)
	if err != nil {
		return
	}
	// 插入新的普通用户数据
	newId, err := models.InsertSimpleUser(newAccount, r.Uid, r.AgentLevel)
	if err != nil {
		return
	}
	// 插入代理关系
	err = models.InsertAgentRelation(r.Uid, newId)
	if err != nil {
		return
	}
	return nil
}
