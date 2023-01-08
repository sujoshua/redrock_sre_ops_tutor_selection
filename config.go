package main

import (
	"github.com/goccy/go-json"
	"log"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Addr       string    `json:"Addr"`      // 服务监听地址
	StartTime  time.Time `json:"StartTime"` // 开始时间
	startTime  int64
	MaxNum     int    `json:"MaxNum"` // 最大允许人数
	TokenKey   string `json:"TokenKey"`
	tokenKey   []byte
	ServerURL  string `json:"ServerURL"` // 当前服务的api
	serverURL  *url.URL
	RosterFile string       `json:"RosterFile"` // 允许获取token的学员名单
	Names      []NameConfig `json:"Names"`      // 导师们的名字
}

type NameConfig struct {
	Name         string `json:"Name"`         // 导师名字
	NickName     string `json:"NickName"`     // 需要用到英文的地方，用来替换姓名
	FeishuURL    string `json:"FeishuURL"`    // 飞书群组url
	SpecialWords string `json:"SpecialWords"` // 成功加入群组后返回想说的特殊的话，可空
}

// GetConfig 解析配置文件
func GetConfig() (config Config, err error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Read config failed, error occur: %s", err.Error())
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Unmarshal json fail, error occur: %s", err.Error())
	}
	config.startTime = config.StartTime.Unix()
	config.tokenKey = []byte(config.TokenKey)

	config.serverURL, err = url.Parse(config.ServerURL)
	if err != nil {
		log.Printf("parse server URL, error occur: %s", err.Error())
	}

	return
}
