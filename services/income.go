package services

import (
	"wr_fenfa/models"
	"wr_fenfa/utils"
)

// 用户id,获取该用户所有上级的id
func SearchAllParent(childId int) (parentId []int, err error) {
	for {
		parentId = append(parentId, childId)
		childId, err = models.GetParent(childId)
		if err != nil {
			if err.Error() == utils.ErrNoRow() {
				err = nil
			}
			return
		}
	}
	return
}

// 根据代理列表和初始代理价格,得到每 个代理的收益
func SearchAllIncome(idListDesc []int, price float64) (result map[int]float64, err error) {
	result = make(map[int]float64)
	for k, id := range idListDesc {
		var income float64
		var x int
		if k != len(idListDesc)-1 {
			x, err = models.QueryParentscale(idListDesc[k+1])
			income = price * float64(100-x) / 100
			price = price - income
			result[id] = income
		} else {
			// 最后一个用户是否是代理
			count, _ := models.SearchIsAgent(id)
			if count > 0 {
				result[id] = price
			} else { // 不是代理的话，将收益转移给上级代理
				result[idListDesc[k-1]] += price
			}
		}
	}
	return
}

// 给每一级代理添加收益记录 并更新总收益
func InsertRegisterIncome(m map[int]float64, cid, pid int, pname string) (err error) {
	err = models.InsertIncomeRecord(m, cid, pid, pname)
	if err != nil {
		return
	}
	err = models.AddUserIncome(m)
	if err != nil {
		return
	}
	return
}
