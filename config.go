package main

import (
	"github.com/goccy/go-json"
	"log"
	"os"
	"time"
)

type Config struct {
	Addr      string    `json:"Addr"`      // 服务监听地址
	StartTime time.Time `json:"StartTime"` // 开始时间
	startTime int64
	MaxNum    int    `json:"MaxNum"` // 最大允许人数
	TokenKey  string `json:"TokenKey"`
	tokenKey  []byte
}

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
	return
}
