package chat

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wsimserver/utils"
)

type Message []byte

//消息内容
type Msg struct {
	Type    string `json:"type"`    //消息类型  txt文本消息（emoji表情） pic图片  video视频 audio语音
	Content string `json:"content"` //消息内容
}

type OutPut struct {
	TargetType string            `json:"target_type"` //user单聊  group群聊
	Target     string            `json:"target"`      //消息接收者
	Msg        Msg               `json:"msg"`         //消息内容
	Ext        map[string]string `json:"ext"`         //扩展属性
	From       string            `json:"from"`        //消息发送者
	Timestamp  int64             `json:"timestamp"`   //消息发送时间
}

type InPut struct {
	TargetType string            `json:"target_type"` //user单聊  group群聊
	Target     string            `json:"target"`      //消息接收者
	Msg        Msg               `json:"msg"`         //消息内容
	Ext        map[string]string `json:"ext"`         //扩展属性
	From       string            `json:"from"`        //消息发送者
}

func (self Message) String() string {
	return string(self)
}
func ParseMessage(msg Message) *InPut {
	var val InPut
	err := json.Unmarshal([]byte(msg), &val)
	if err != nil {
		return nil
	}
	if val.TargetType != "user" && val.TargetType != "group" {
		return nil
	}
	if len(val.Target) == 0 {
		return nil
	}
	if len(val.From) == 0 {
		return nil
	}
	if val.Msg.Type != "txt" {
		return nil
	}
	return &val
}

func NewOutput(val *InPut) *OutPut {
	return &OutPut{
		TargetType: val.TargetType,
		Target:     val.Target,
		Msg:        val.Msg,
		Ext:        val.Ext,
		From:       val.From,
		Timestamp:  time.Now().Unix(),
	}
}

//struct to json string
func (temp *OutPut) String() string {
	b, err := json.Marshal(temp)
	if err != nil {
		return ""
	}
	return string(b)
}

//struct to json []byte
func (temp *OutPut) Bytes() Message {
	b, err := json.Marshal(temp)
	if err != nil {
		return nil
	}
	return b
}

func WhetherLogin(msg Message) (bool, bool, string) {
	var (
		action, login bool
		userid        string
	)
	temp := string(msg)
	if strings.Contains(temp, "token=Bearer") {
		tokenTemp := strings.Split(temp, "=")
		if len(tokenTemp) > 1 {
			tokenString := tokenTemp[1]
			user, err := utils.ParseJWTokenUserId(tokenString)
			if err == nil {
				action = true
				login = true
				userid = fmt.Sprintf("%d", user)
			} else {
				action = true
				login = false
				userid = "0"
			}
		} else {
			action = true
			login = false
			userid = "0"
		}
	} else {
		action = true
		login = false
		userid = "0"
	}
	return action, login, userid
}
