/**
* @Author:changjiang
* @Description:
* @File:http
* @Version: 1.0.0
* @Date 2020/6/10 7:05 下午
 */
package export

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/dashboard-api-main/web/pojos"
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
	pojos.HttpRequestContent
}

func NewHttpRequest(context *base.Context) (r *HttpRequest) {
	r = &HttpRequest{
		Context: context,
	}
	return
}

func (r *HttpRequest) sendPost2() {
	apiUrl := "http://127.0.0.1"
	resource := "/tpost"
	data := url.Values{}
	data.Set("name", "rnben")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "http://127.0.0.1/tpost"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println("post send success")

}

// 发送HTTp请求接口
func (r *HttpRequest) Send(request *RequestObject) (res *[]byte, err error) {

	r.Context.Log.Infoln("request:", fmt.Sprintf("%v", request))

	res = &[]byte{}
	var urlVal = request.Uri

	argString := make([]string, 0)
	for key, value := range request.Query {
		argString = append(argString, fmt.Sprintf("%s=", key)+r.getParseParam(fmt.Sprintf("%v", value)))
	}
	tp := strings.Join(argString, "&")
	urlArr := strings.Split(urlVal, "?")
	len := len(urlArr)
	if len == 1 {
		urlVal = urlArr[0] + "?" + tp
	} else if len == 2 {
		urlVal = urlArr[0] + "?" + r.getParseParam(urlArr[1]) + tp
	} else {
		r.Context.Log.Errorln(fmt.Sprintf("发送HTTP请求 参数异常:%s", urlVal))
	}
	client := &http.Client{}
	var req *http.Request
	if request.Body == nil {
		r.Context.Log.Infoln("request:", fmt.Sprintf("%s", urlVal))
		req, _ = http.NewRequest(request.RequestMethod, urlVal, nil)
	} else {
		r.Context.Log.Infoln("request:", fmt.Sprintf("%s", urlVal))
		req, _ = http.NewRequest(request.RequestMethod, urlVal, request.Body)
	}
	for headerKey, v := range request.HttpHeader {
		// 前端发送请求 "X-"开头的的header将会自动上下文传递
		if strings.HasPrefix(strings.ToUpper(headerKey), "X-") {
			req.Header.Add(headerKey, strings.Join(v, ""))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		r.Context.Log.Errorln("request:", fmt.Sprintf("%v", request), "content:", err.Error())
		return
	}
	defer resp.Body.Close()
	*res, _ = ioutil.ReadAll(resp.Body)
	return
}

// 将get请求的参数进行转义
func (r *HttpRequest) getParseParam(param string) string {
	return url.PathEscape(param)
}
