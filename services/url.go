package services

import (
	"wr_fenfa/utils"
	"zcm_tools/http"
	"encoding/json"
)

type Url struct {
	Result   bool   `json:"result"`
	UrlShort string `json:"url_short"`
	UrlLong  string `json:"url_long"`
	Type     int    `json:"type"`
}
type LongUrl struct {
	Urls []Url `json:"urls"`
}

//长链接到短链接
func LongToSort(urlLong string) (newUrl string, err error) {
	url := utils.LANGTOSHORTURL + "?access_token=" + utils.ACCESSTOKEN + "&url_long=" + urlLong
	b, err := http.Get(url)
	if err != nil {
		return "", err
	} else {
		var m LongUrl
		if err := json.Unmarshal(b, &m); err == nil {
			if len(m.Urls) > 0 {
				return m.Urls[0].UrlShort, err
			}
		}
		return "", err
	}
}
