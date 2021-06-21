package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

var allowMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"HEAD":    true,
	"PUT":     true,
	"PATCH":   true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
}

type HttpClient struct {
	ConnTimeOut  int
	ReadTimeOut  int
	WriteTimeOut int

	url    string
	method string
	header map[string]interface{}
	params map[string]interface{}
	body   string
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		//ConnTimeOut:  1000, //millisecond
		//ReadTimeOut:  2000,
		//WriteTimeOut: 1000,

		url: url,
	}
}

func (i *HttpClient) T(c, r, w int) *HttpClient {
	i.ConnTimeOut = c
	i.ReadTimeOut = r
	i.WriteTimeOut = w
	return i
}

func (i *HttpClient) M(m string) *HttpClient {
	i.method = m
	return i
}

func (i *HttpClient) H(k string, v interface{}) *HttpClient {
	if i.header == nil {
		i.header = make(map[string]interface{})
	}
	i.header[k] = v
	return i
}

func (i *HttpClient) P(k string, v interface{}) *HttpClient {
	if i.params == nil {
		i.params = make(map[string]interface{})
	}
	i.params[k] = v
	return i
}

func (i *HttpClient) B(b string) *HttpClient {
	i.body = b
	return i
}

type Json struct {
	Data []byte
}

func (j *Json) Get(keys ...string) (value []byte, dataType jsonparser.ValueType, err error) {
	v, t, _, err := jsonparser.Get(j.Data, keys...)
	return v, t, err
}

func (j *Json) Unmarshal(val interface{}) (err error) {
	err = json.Unmarshal(j.Data, val)
	if err != nil {
		return err
	}
	return nil
}

func (h *HttpClient) Do() (retCode int, json *Json, err error) {
	if _, ok := allowMethods[h.method]; !ok {
		return 0, nil, errors.New("Current method is not allowed.")
	}

	if len(h.params) > 0 {
		if !strings.ContainsRune(h.url, '?') {
			h.url = fmt.Sprintf("%s?", h.url)
		} else {
			h.url = fmt.Sprintf("%s&", h.url)
		}

		for k, v := range h.params {
			h.url = fmt.Sprintf("%s%s=%v&", h.url, k, v)
		}
		h.url = h.url[:len(h.url)-1]
	}

	client := &http.Client{
		//Transport: &http.Transport{
		//	ResponseHeaderTimeout: time.Millisecond * time.Duration(h.ReadTimeOut),
		//},
		Timeout: time.Millisecond * time.Duration(h.ConnTimeOut+h.ReadTimeOut+h.WriteTimeOut),
	}

	req, err := http.NewRequest(h.method, h.url, strings.NewReader(h.body))
	if err != nil {
		return 0, nil, err
	}

	for k, v := range h.header {
		req.Header.Add(k, v.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	json = &Json{
		Data: []byte{},
	}
	json.Data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, json, nil
}
