package services

import (
	"net/url"
	"time"
	"wr_fenfa/utils"
	"encoding/xml"
	"zcm_tools/http"
)

/*
	第三方短信发送
*/
func SendXinShimessage(content, mobile string) (xsSmsResult XsSmsResult, err error) {
	timeNow := time.Now().Format("20060102150405")
	sign := utils.XS_NAME + utils.XS_PASSWORD + timeNow
	sign = utils.MD5(sign)
	params := url.Values{}
	params.Set("userid", utils.XS_USERID)
	params.Set("timestamp", timeNow)
	params.Set("sign", sign)
	params.Set("mobile", mobile)
	params.Set("content", content)
	params.Set("sendTime", "")
	params.Set("action", "send")
	params.Set("extno", "")
	s, err := http.Post(utils.XS_URL, params.Encode())
	if err != nil {
		return
	}
	err = xml.Unmarshal(s, &xsSmsResult)
	if err != nil {
		return
	}
	return
}
