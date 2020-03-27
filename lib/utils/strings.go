/**
* @Author:changjiang
* @Description:
* @File:strings
* @Version: 1.0.0
* @Date 2020/3/19 11:45 下午
 */
package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	c := md5.New()
	c.Write([]byte(s))
	cipherStr := c.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
