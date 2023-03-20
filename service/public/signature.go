package public

import (
	"crypto/sha1"
	"sort"
	"strings"
	"fmt"
	"log"
)

func CheckSignature(signature string, timestamp string, nonce string,token string)(bool) {
	//将token、timestamp、nonce三个参数放到一个数组中
	strs:=[]string{token,timestamp,nonce}
	//将数组中的元素按照字典序排序
	strs=sort.StringSlice(strs)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	joinStr:=strings.Join(strs,"")
	//生成签名
	sha1:=sha1.Sum([]byte(joinStr))
	//将签名转为字符串
	sha1Str:=fmt.Sprintf("%x",sha1)
	//将sha1加密后的字符串与signature进行对比
	if sha1Str==signature {
		return true
	}
	log.Printf("checkSignature failed,sha1 %s,signature %s \n",sha1Str,signature)
	return false
}