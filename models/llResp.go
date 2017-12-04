package models

// 连连回调参数
type RechargeFeedback struct {
	Ret_code    string
	Ret_msg     string
	Sign_type   string
	Sign        string
	Oid_partner string
	No_order    string
	Dt_order    string
	Money_order string
	Oid_paybill string
	Result_pay  string
	Settle_date string
	Info_order  string
}

//连连放款订单查询
type LLTradeQueryResponse struct {
	Ret_code    string
	Ret_msg     string
	Sign_type   string
	Sign        string
	Oid_partner string
	Result_pay  string
	Dt_order    string
	No_order    string
	Oid_paybill string
	Meney_order string
	Settle_date string
	Info_order  string
}

//连连银行卡卡bin查询接口 返回数据
type LLBandCardBinResponse struct {
	Ret_code  string
	Ret_msg   string
	Sign_type string
	Sign      string
	Bank_code string
	Bank_name string
	Card_type string
}

type LLTradeApiResponse struct {
	Ret_code     string
	Ret_msg      string
	Sign_type    string
	Sign         string
	Oid_partner  string
	No_order     string
	Oid_paybill  string
	Confirm_code string
	Token        string
}

type LLTradeResponse struct {
	Ret_code     string
	Ret_msg      string
	Sign_type    string
	Sign         string
	Oid_partner  string
	No_order     string
	Oid_paybill  string
	Confirm_code string
}