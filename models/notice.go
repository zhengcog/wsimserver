package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const NOTICEURL = "http://127.0.0.1/umeng/send?"

func NoticeUMeng(user string, message string) {
	query := url.Values{"users": {user}, "msgtype": {"chat"}, "message": {message}}
	url := query.Encode()
	res, err := http.Get(NOTICEURL + url)
	if err != nil {
		return
	}
	detail, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	_ = detail
}

func ParseString(msg string) (string, string) {
	i := strings.LastIndex(msg, "_")
	if i < 0 {
		return "", ""
	}
	i += 1
	uDec, err := base64.StdEncoding.DecodeString(msg[i:])
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	all := make(map[string]interface{})
	json.Unmarshal(uDec, &all)
	result := all["conversation"]
	conversation := result.(string)
	var detail string
	switch {
	case strings.Contains(msg, "text"):
		{
			i := strings.Index(msg, "_")
			i += 1
			msg = msg[i:]
			j := strings.Index(msg, "_")
			j += 1
			msg = msg[j:]
			k := strings.Index(msg, "_")
			detail = ":" + msg[:k]
		}
	case strings.Contains(msg, "emotion"):
		{
			detail = ":[表情]"
		}
	case strings.Contains(msg, "picture"):
		{
			detail = ":[图片]"
		}
	case strings.Contains(msg, "video"):
		{
			detail = ":[视频]"
		}
	case strings.Contains(msg, "audio"):
		{
			detail = ":[语音]"
		}
	}

	return conversation, detail
}

func UMengNotice(user string, msg string) {
	message, detail := ParseString(msg)
	NoticeUMeng(user, message+detail)
}
