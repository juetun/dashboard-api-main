/**
* @Author:changjiang
* @Description:
* @File:http
* @Version: 1.0.0
* @Date 2020/6/10 7:05 下午
 */
package export

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/wrappers"
)

type HttpRequest struct {
	Context *base.Context
}
type RequestObject struct {
	Uri           string    `json:"uri"`
	ContentType   string    `json:"content_type"`
	Body          io.Reader `json:"body"`
	Query         map[string]interface{}
	RequestMethod string `json:"request_method"`
	ReSendTimes   int // 如果获取数据失败 尝试重新请求的次数 默认3次
	wrappers.HttpRequestContent
}

func NewHttpRequest(context *base.Context) (r *HttpRequest) {
	r = &HttpRequest{
		Context: context,
	}
	return
}

// func (r *HttpRequest) sendPost2() {
// 	apiUrl := "http://127.0.0.1"
// 	resource := "/tpost"
// 	data := url.Values{}
// 	data.Set("name", "rnben")
//
// 	u, _ := url.ParseRequestURI(apiUrl)
// 	u.Path = resource
// 	urlStr := u.String() // "http://127.0.0.1/tpost"
//
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
//
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer resp.Body.Close()
// 	fmt.Println("post send success")
//
// }

func (r *HttpRequest) preReadySend(request *RequestObject, client *http.Client, req *http.Request, res *[]byte, timeStep int) {
	timeStep--
	r.Context.Info(map[string]interface{}{
		"request:": fmt.Sprintf("%v", request),
		"desc:":    fmt.Sprintf("尝试第%d次(共%d次)获取数据", request.ReSendTimes-timeStep, request.ReSendTimes),
	})
	if timeStep < 1 {
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		r.Context.Error(map[string]interface{}{
			"request:": fmt.Sprintf("%v", request), "content:": err.Error(),
		})

		// 等待一定时间后重新发送获取数据的指令 4的倍数秒 第一次:4秒 ,第二次:8秒, 第三次:12秒
		time.Sleep(time.Duration((request.ReSendTimes-timeStep)*4) * time.Second)
		r.preReadySend(request, client, req, res, timeStep)
		return
	}
	*res, _ = ioutil.ReadAll(resp.Body)
}

func (r *HttpRequest) orgParams(request *RequestObject) (logData map[string]interface{}, urlVal string) {
	logData = map[string]interface{}{
		"request content:": fmt.Sprintf("%v", request),
		"Method":           request.RequestMethod,
	}
	urlVal = request.Uri

	argString := make([]string, 0)
	for key, value := range request.Query {
		argString = append(argString, fmt.Sprintf("%s=", key)+r.getParseParam(fmt.Sprintf("%v", value)))
	}
	tp := strings.Join(argString, "&")
	urlArr := strings.Split(urlVal, "?")
	len := len(urlArr)
	if len == 1 {
		urlVal = urlArr[0] + "?" + tp
		return
	}
	if len == 2 {
		urlVal = urlArr[0] + "?" + r.getParseParam(urlArr[1]) + tp
		return
	}
	logData["request content:"] = fmt.Sprintf("发送HTTP请求 参数异常:%s", urlVal)
	r.Context.Info(logData)
	return
}

// 发送HTTp请求接口
// @timeStep
func (r *HttpRequest) Send(request *RequestObject) (res *[]byte, err error) {
	res = &[]byte{}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, time.Second*2) // 设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2)) // 设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}
	logData, urlVal := r.orgParams(request)

	var req *http.Request
	if request.Body == nil {
		logData["RequestUrl:"] = fmt.Sprintf("%s", urlVal)
		req, _ = http.NewRequest(request.RequestMethod, urlVal, nil)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		s := buf.String()
		logData["RequestUrl:"] = fmt.Sprintf("%s", urlVal)
		logData["Body:"] = s
		req, _ = http.NewRequest(request.RequestMethod, urlVal, request.Body)
	}
	for headerKey, v := range request.HttpHeader {
		// 前端发送请求 "X-"开头的的header将会自动上下文传递
		if strings.HasPrefix(strings.ToUpper(headerKey), "X-") {
			vc := strings.Join(v, "")
			logData[headerKey] = vc
			req.Header.Add(headerKey, vc)
		}
	}
	r.Context.Info(logData)
	// 尝试获取数据，如果获取失败，则多次尝试
	r.preReadySend(request, client, req, res, request.ReSendTimes)
	return
}

// 将get请求的参数进行转义
func (r *HttpRequest) getParseParam(param string) string {
	return url.PathEscape(param)
}
