package models

import (
	"log"
	"wsimserver/gosexy/redis"
	"wsimserver/utils"
)

const REDISKEY = "offline_msg_"

//友盟推送消息
func PushOfflineMsg(user string, msg string) {
	var (
		client *redis.Client
		ok     bool
	)
	client, ok = utils.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	redisKey := REDISKEY + user
	client.RPush(redisKey, msg)
	client.Close()
	//go UMengNotice(user, msg)
}

func SendoutOfflineMsg(user string) []string {
	var (
		client *redis.Client
		ok     bool
		null   []string
	)
	client, ok = utils.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return null
	}
	redisKey := REDISKEY + user
	msg, _ := client.LRange(redisKey, 0, -1)
	client.Del(redisKey)
	client.Close()
	return msg
}
