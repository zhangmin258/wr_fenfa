package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// 产品信息模块
type Product struct {
	Id         int
	Name       string    // 产品名称
	ProductId  int       // 产品id
	AgentPrice float64   // 产品价格
	CpaDefine  string    // 产品有效性定义
	IsUse      int       // 产品状态 0，使用 1，不使用
	CreateTime time.Time // 创建时间
}

// 修改产品信息
type ProductEdit struct {
	ProductId  int     // 产品id
	AgentPrice float64 // 产品价格
}

// 获取产品列表
func GetProductList(condition string, params []string, begin, size int) (product []Product, err error) {
	o := orm.NewOrm()
	o.Using("wr")
	sql := `SELECT id AS product_id,name,is_use,agent_price,"注册" as cpa_define,create_time FROM product WHERE cooperation_type=1 AND loan_product_type=0 AND id!=0 `
	sql += condition
	sql += ` ORDER BY create_time DESC LIMIT ?,? `
	_, err = o.Raw(sql, params, begin, size).QueryRows(&product)
	return
}

// 获取产品列表
func GetProductCount(condition string, params []string) (count int, err error) {
	o := orm.NewOrm()
	o.Using("wr")
	sql := `SELECT COUNT(1) FROM product WHERE cooperation_type=1 AND loan_product_type=0 AND id!=0 `
	sql += condition
	err = o.Raw(sql, params).QueryRow(&count)
	return
}

// 编辑产品
func UpdateProduct(productEdit ProductEdit) error {
	o := orm.NewOrm()
	o.Using("wr")
	sql := `UPDATE product SET agent_price = ? WHERE id = ? `
	_, err := o.Raw(sql, productEdit.AgentPrice, productEdit.ProductId).Exec()
	return err
}

// 根据id获取产品代理价
func GetProductAgentPrice(pid int) (p Product, err error) {
	o := orm.NewOrm()
	o.Using("wr")
	sql := `SELECT id,name,agent_price FROM product WHERE id=?`
	err = o.Raw(sql, pid).QueryRow(&p)
	return
}
