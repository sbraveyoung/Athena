package httpClient

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

type resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

//func TestHttpClient(t *testing.T) {
//
//	//	url := "http://a23-scalper.vipkid-qa.com.cn/search/onlineClass/9998824751"
//
//	url := "http://10.106.128.62/api/vps/service/openclass/s2/v1/courseware/getcourseinf?vcwId=cwf17c5278f6cc5c73b9dcc170552c5526"
//	ret, json, err := NewHttpClient(url).M("GET").T(3000, 3000, 3000).Do()
//	if err == nil && ret == http.StatusOK && json != nil {
//		//case 1
//		value, dataType, err := json.Get("data")
//		fmt.Println(string(value), dataType, err)
//
//		//case 2
//		var r resp
//		if json.Unmarshal(&r) != nil {
//			t.Error("test fail. json.Unmarshal fail,err:", err)
//		}
//		fmt.Println(r.Data)
//	} else {
//		t.Error("test fail. url", url, "ret:", ret, "json:", json, "err:", err.Error())
//	}
//
//}

func TestHttpClient(t *testing.T) {

	url := "http://10.106.128.62/api/vps/service/openclass/s2/v1/courseware/getcourseinf?vcwId=cwf17c5278f6cc5c73b9dcc170552c5526"
	count := int64(0)
	for {
		go func() {
			resp, err := http.Get(url)
			if err == nil {
				defer resp.Body.Close()
				//ret, json, err := NewHttpClient(url).M("GET").T(3000, 3000, 3000).Do()
				//if err == nil && ret == http.StatusOK && json != nil {
				//json.Get("data")
				//value, dataType, err := json.Get("data")
				//fmt.Println(string(value), dataType, err)
				//value, _, _ := json.Get("data")
				//fmt.Println(string(value))
			} else {
				fmt.Println("Benchmark error:", err)
				os.Exit(1)
			}
		}()
		count++
		fmt.Println(count)
		time.Sleep(100)
	}
}
