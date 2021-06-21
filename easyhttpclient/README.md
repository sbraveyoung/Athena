
**httpc** is a http client with go that easy to use.

### Usage

```go
	url := "http://www.baidu.com"
	ret, jdata, err := NewHttpClient(url).M("GET").T(1000,1000,1000).H("Auth", "111").P("x", "222").B("test").Do()
	if err == nil && ret == http.StatusOK && jdata != nil {
		//case 1
		value, dataType, err := json.Get("data")
		fmt.Println(string(value), dataType, err)

		//case 2
		var r resp
		if json.Unmarshal(&r) != nil {
			t.Error("test fail. json.Unmarshal fail,err:", err)
		}
		fmt.Println(r.Data)
	} else {
		t.Error("test fail. url", url, "ret:", ret, "jdata:", jdata, "err:", err.Error())
	}
```

**Please Note:If you get a nil value from `json.Get()` function, check the response of http first, maybe the returned data does not confirm to the json format specifiction. For example:**

> {"code":"200","msg":"Success: ","data":"{\"questionCount\":4,\"videoDuration\":829240,\"durationMax\":1066530,\"durationMin\":1061602,\"videoStatus\":0,\"videoUrl\":\"\",\"eventUrl\":\"https://openclass-cdn.vipkid.com.cn/gateway/beta/text/txtf182ad040aa38827b2e5dde1fc296630.json\"}"}

The type of data's value is string,instead of object.

The allowed [json object](http://json.org) types are those:

* String
* Number
* Object
* Array
* Boolean
* Null
