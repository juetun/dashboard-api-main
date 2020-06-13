/**
* @Author:changjiang
* @Description:
* @File:strings.
* @Version: 1.0.0
* @Date 2020/6/13 1:07 下午
 */
package utils

import (
	"fmt"
	"time"
)

func ShowDateTime(timeValue time.Time) (res string) {
	second := time.Now().Unix() - timeValue.Unix()
	day := second / 86400
	if day > 1 {
		res = fmt.Sprintf("%d天前", day)
		return
	}
	daySecond := second % 86400
	dayHour := daySecond / 3600
	if dayHour > 1 {
		res = fmt.Sprintf("%d小时前", dayHour)
		return
	}
	daySecond = daySecond % 3600
	dayMinute := daySecond / 60
	if dayMinute > 1 {
		res = fmt.Sprintf("%d分钟前", dayMinute)
		return
	}
	res = "刚刚"
	return
}
